package ssh

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

/**
	Copy a file or a directory from the source to the remote destination
 */
func (client *Client) Copy(srcPath, destPath string) error {
	handle, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	// Load file stats
	stats, err := handle.Stat()
	if err != nil {
		return err
	}
	isDir := stats.IsDir()
	handle.Close()

	if isDir {
		return client.CopyDir(srcPath, destPath)
	}
	return client.CopyFile(srcPath, destPath)
}

/**
	Copy a file from the source to the destination
 */
func (client *Client) CopyFile(srcPath, destination string) error {
	// start SSH connection
	session, err := client.StartSession(false, false)
	if err != nil {
		return fmt.Errorf("Unable to contact server[%s]: %s", client.Config.Host, err)
	}
	defer session.Close()

	// open an input stream to the server
	dest, _ := session.StdinPipe()
	defer func(){
		if dest!=nil {
			dest.Close()
		}
	}()
	destPath := path.Base(destination)

	dir := path.Dir(destination)
	dir = filepath.ToSlash(dir)

	// start receiving the file on the server using scp but don't wait for the command to finish
	cmd := fmt.Sprintf("scp -vt %s", dir)
	if err := session.Start(cmd); err != nil {
		session.Close()
		return err
	}

	// send the file over the wire
	if err = transferFile(srcPath, destPath, dest); err != nil {
		return err
	}
	dest.Close(); dest = nil

	// wait until the command has finished and see if there are any errors
	err = session.Wait()
	if err != nil {
		return err
	}

	return nil
}

/**
	Open a file and transfer it to the destination
 */
func transferFile(srcPath string, destPath string, dest io.Writer) error {
	// Open file for reading
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Load file stats
	stats, err := src.Stat()
	if err != nil {
		return err
	}
	size := stats.Size()
	mode := stats.Mode().Perm()

	// Send content through the connection
	if err = scpTransferFile(destPath, mode, size, src, dest); err != nil {
		return err
	}

	return nil
}

/**
	scpTransferFile sends the contents of the source stream to the destination stream using SCP protocol
 */
func scpTransferFile(path string, mode os.FileMode, size int64, src io.Reader, dest io.Writer) error {

	// send the location where it should be saved, along with the size and mode
	fmt.Fprintf(dest, "C%#o %d %s\n", mode, size, path)

	// then send the contents of the file
	if _, err := io.Copy(dest, src); err != nil {
		return err
	}

	// complete the file transfer with a null value
	fmt.Fprint(dest, "\x00")

	return nil
}