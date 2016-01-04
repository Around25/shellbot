package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

/**
	Copy a directory from the source to the destination
 */
func (client *Client) CopyDir(srcPath, destination string) error {
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

	// start receiving the file on the server using scp but don't wait for the command to finish
	destination = filepath.ToSlash(destination)
	cmd := fmt.Sprintf("scp -rvt %s", destination)
	if err := session.Start(cmd); err != nil {
		session.Close()
		return err
	}

	// send the directory over the wire
	if err = transferDir(srcPath, dest); err != nil {
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
	transferDir first creates the folder on the server and them transfers all it's contents on the server
 */
func transferDir(srcPath string, dest io.Writer) error {
	// open the provided source directory
	handle, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		if handle != nil {
			handle.Close()
		}
	}()
	// Load dir stats
	stats, err := handle.Stat()
	if err != nil {
		return err
	}
	name := stats.Name()
	mode := stats.Mode().Perm()
	handle.Close(); handle = nil

	// transfer the current folder first
	err = scpTransferDir(name, mode, dest, func() error {
		return transferDirContents(srcPath, dest)
	})
	if err != nil {
		return err
	}
	return nil
}

/**
	Transfer directory recursively to the destination
 */
func transferDirContents(srcPath string, dest io.Writer) error {
	// open the provided source directory
	handle, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer handle.Close()

	// read all contents
	entries, err := handle.Readdir(-1)
	if err != nil {
		return err
	}

	// traverse each item and transfer the right data over SCP for each one
	for _, fileInfo := range entries {
		fullFilePath := filepath.Join(srcPath, fileInfo.Name())

		if !fileInfo.IsDir() {
			transferFile(fullFilePath, fileInfo.Name(), dest)
			continue
		}

		err = scpTransferDir(fileInfo.Name(), fileInfo.Mode().Perm(), dest, func() error {
			return transferDirContents(fullFilePath, dest)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

/**
	scpTransferDir implements the SCP protocol for creating a directory at the destination
 */
func scpTransferDir(path string, mode os.FileMode, dest io.Writer, processFiles func() error) error {

	// send the location where a folder should be created
	fmt.Fprintf(dest, "D%#o %d %s\n", mode, 0, path)

	// transfer the files contained in the folder
	if err := processFiles(); err != nil {
		return err
	}

	// Send the close folder message
	fmt.Fprint(dest, "E\n")
	return nil
}