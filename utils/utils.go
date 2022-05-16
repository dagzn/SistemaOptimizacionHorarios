package utils

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	obj "proyecto-horarios/objetos"
)

func LeerArchivo(archivo string) ([]byte, error) {
	data, err := ioutil.ReadFile(archivo)
	return data, err
}

func DeserializarEntradaHorario(data []byte) (*obj.Entrada_horario, error) {
	var h *obj.Entrada_horario
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func SerializarSalidaHorario(h *obj.Salida_horario) ([]byte, error) {
	data, err := json.MarshalIndent(h, "", " ")
	if err != nil {
		return nil, err
	}
	return data, err
}

func DeserializarEntradaValidacion(data []byte) (*obj.Entrada_validacion, error) {
	var h *obj.Entrada_validacion
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func SerializarSalidaValidacion(h *obj.Salida_validacion) ([]byte, error) {
	data, err := json.MarshalIndent(h, "", " ")
	if err != nil {
		return nil, err
	}
	return data, err
}

func DeserializarEntradaExportacion(data []byte) (*obj.Entrada_exportacion, error) {
	var h *obj.Entrada_exportacion
	err := json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func SerializarSalidaExportacionFallida(h *obj.Salida_exportacion_fallida) ([]byte, error) {
	data, err := json.MarshalIndent(h, "", " ")
	if err != nil {
		return nil, err
	}
	return data, err
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
