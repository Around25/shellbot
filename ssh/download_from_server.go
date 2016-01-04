package ssh

import (
	"io"
	"os"
	"fmt"
	"path"
	"bufio"
	"strconv"
	"github.com/Around25/shellbot/logger"
	"path/filepath"
)

/**
	Download a file or directory from the server
 */
func (client *Client) Download(srcPath, destination string) error {
	// start SSH connection
	session, err := client.StartSession(false, false)
	if err != nil {
		return fmt.Errorf("Unable to contact server[%s]: %s", client.Config.Host, err)
	}
	defer session.Close()

	// open an input stream to the server
	dest, _ := session.StdinPipe()
	defer func() {
		if dest != nil {
			dest.Close()
		}
	}()

	sourceStream, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	source := bufio.NewReader(sourceStream)

	// start receiving the data from scp
	cmd := fmt.Sprintf("scp -vrf %s", strconv.Quote(srcPath))
	if err := session.Start(cmd); err != nil {
		session.Close()
		return err
	}

	// send the directory over the wire
	if err = scpReceive(source, dest, destination); err != nil {
		return err
	}
	dest.Close(); dest = nil

	// wait until the command has finished and see if there are any errors
	err = session.Wait()
	if err != nil {
		logger.Warning("SCP command failed")
		return err
	}

	return nil
}

/**
	Read the next message from the scp stream
 */
func readHeader(source *bufio.Reader) (string, error) {
	header, err := source.ReadString('\n')
	if err == io.EOF {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return header, nil
}

/**
	Extract file info from the header
 */
func readFileInfo(header string) (os.FileMode, int64, string, error) {
	var mode os.FileMode
	var size int64
	var name string
	var messageType string

	n, err := fmt.Sscanf(header, "%1s%o %d %s\n", &messageType, &mode, &size, &name)
	if err != nil || n != 4 {
		return mode, size, name, fmt.Errorf("Invalid response from server: %s %s", header, err)
	}
	return mode, size, name, nil
}

/**
	Receive all SCP data from the source and save it in the destination
	Send confirmation messages using the reply stream
 */
func scpReceive(source *bufio.Reader, reply io.Writer, dest string) error {
	// confirm communication chanel
	fmt.Fprint(reply, "\x00");

	// start processing messages from the server
	if err := scpProcessMessage(source, reply, dest); err != nil {
		return err;
	}

	return nil
}

/**
	Receive a single file from the server
 */
func scpReceiveFile(header string, source *bufio.Reader, reply io.Writer, dest string) error {
	// read the file info from the message header
	mode, size, name, err := readFileInfo(header)
	if err != nil {
		return err
	}

	// confirm receiving the properties of the file
	fmt.Fprint(reply, "\x00")
	var filename string
	if filepath.Ext(dest) != "" && dest[len(dest) - 1] != byte(filepath.Separator) {
		filename = dest
	} else {
		filename = path.Join(dest, name)
	}

	// open the file in which to save the data
	f, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	// don't forget to close the file at the end
	defer f.Close()

	// copy all the data from the server to the file
	if _, err := io.CopyN(f, source, size); err != nil {
		return err
	}

	// confirm receiving the contents of the file
	fmt.Fprint(reply, "\x00")

	emptyChar, err := source.ReadByte()
	if err != nil || emptyChar != '\x00' {
		return fmt.Errorf("Invalid char at the end of file")
	}

	return nil
}

/**
	Process a directory scp message and create it in the destination
 */
func scpReceiveDir(header string, source *bufio.Reader, reply io.Writer, dest string) (string, error) {
	// read the folder details from the message header
	mode, _, name, err := readFileInfo(header)
	if err != nil {
		return dest, err
	}

	// get the new directory path that should be created
	dest = path.Join(dest, name)
	// create the full path
	err = os.MkdirAll(dest, mode)

	// if an error is found send proper reply to the server
	if err != nil {
		fmt.Fprint(reply, "\x01")
		return dest, err
	}

	// confirm receiving the properties of the folder
	fmt.Fprint(reply, "\x00")

	return dest, nil
}

/**
	Process SCP messages one at a time until the stream is over
 */
func scpProcessMessage(source *bufio.Reader, reply io.Writer, dest string) error {
	header, err := readHeader(source)
	if err != nil {
		return err
	}

	if len(header) == 0 {
		return nil // nothing left to receive
	}

	switch (header[0]) {
	default:
		// handle bad data received
		return fmt.Errorf("Invalid message received '%s'", header)
	case '\x01', '\x02':
		// handle any errors in communication
		return fmt.Errorf("%s", header[1:len(header)])
	case 'C':
		// receive a file and save it in the given folder
		err = scpReceiveFile(header, source, reply, dest);
		if err != nil {
			return err
		}
	case 'D':
		// receive a directory name and create it
		dest, err = scpReceiveDir(header, source, reply, dest);
		if err != nil {
			return err
		}
	case 'E':
		// close a directory and move to the parent one
		dest = path.Dir(dest)
		fmt.Fprint(reply, "\x00")
	case 'T':
	// ignore time messages
	}
	// continue processing messages until the buffer is empty
	return scpProcessMessage(source, reply, dest)
}