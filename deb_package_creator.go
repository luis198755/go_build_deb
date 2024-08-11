package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Bienvenido al creador de paquetes .deb")

	// Recopilar información del paquete
	packageName := promptUser("Nombre del paquete")
	version := promptUser("Versión")
	architecture := promptUser("Arquitectura (e.g., all, amd64)")
	maintainer := promptUser("Mantenedor (Nombre <email>)")
	shortDescription := promptUser("Descripción corta")
	longDescription := promptUser("Descripción larga (use '.' para terminar)")
	executablePath := promptUser("Ruta del archivo ejecutable a incluir")

	// Crear un archivo temporal para el script de Bash
	tmpfile, err := os.CreateTemp("", "create_deb_*.sh")
	if err != nil {
		fmt.Println("Error al crear archivo temporal:", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	// Escribir el script de Bash en el archivo temporal
	bashScript := `#!/bin/bash
mkdir -p "$1/DEBIAN" "$1/usr/local/bin" || exit 1
cat > "$1/DEBIAN/control" << EOF
Package: $2
Version: $3
Architecture: $4
Maintainer: $5
Description: $6
 $7
EOF
cp "$8" "$1/usr/local/bin/" || exit 1
sudo chown -R root:root "$1"
sudo chmod -R 755 "$1"
dpkg-deb --build "$1" || exit 1
mv "${1}.deb" "${2}_${3}_${4}.deb" || exit 1
echo "Paquete .deb creado con éxito: ${2}_${3}_${4}.deb"
`
	if _, err := tmpfile.Write([]byte(bashScript)); err != nil {
		fmt.Println("Error al escribir en archivo temporal:", err)
		return
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Println("Error al cerrar archivo temporal:", err)
		return
	}

	// Hacer el script ejecutable
	if err := os.Chmod(tmpfile.Name(), 0755); err != nil {
		fmt.Println("Error al cambiar permisos del archivo temporal:", err)
		return
	}

	// Ejecutar el script de Bash
	cmd := exec.Command("bash", tmpfile.Name(), packageName, packageName, version, architecture, maintainer, shortDescription, longDescription, executablePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error al ejecutar el script:", err)
		return
	}
}

func promptUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}