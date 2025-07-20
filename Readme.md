# Corrector de Scripts de SQ

Esta es una herramienta simple para parchear archivos `.mq5`, reemplazando una línea de código específica con una versión corregida. Se puede ejecutar desde la línea de comandos o a través de una interfaz gráfica de usuario (GUI) instalando accesos directos en el menú contextual.

## El Problema

Los scripts `.mq5` originales contienen un cálculo incorrecto para `PointValue`. Esta herramienta reemplaza la línea defectuosa con una versión corregida que utiliza `SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_CONTRACT_SIZE)`.

## Características

- **Interfaz de Línea de Comandos (CLI):** Procesa archivos directamente desde la terminal.
- **Interfaz Gráfica de Usuario (GUI):** Una ventana simple que muestra el estado de cada archivo que se está procesando.
- **Integración con Menú Contextual:** Haz clic derecho en archivos `.mq5`, carpetas o en el fondo de una carpeta en el Explorador de Windows para parchear los archivos.

## Instalación

1.  Asegúrate de tener PowerShell y Go instalados en tu sistema.
2.  Ejecuta el script `build.ps1` para compilar la aplicación. Esto creará un ejecutable `fix-SQ-scripts.exe` en el directorio del proyecto.
3.  Haz clic derecho en el script `install.ps1` y selecciona "Ejecutar con PowerShell". Esto instalará los accesos directos del menú contextual.

## Cómo Usar

### Modo GUI (a través del Menú Contextual)

-   **Para un solo archivo `.mq5`:** Haz clic derecho en el archivo y selecciona "Fix MQ5 Scripts".
-   **Para una carpeta:** Haz clic derecho en la carpeta y selecciona "Fix MQ5 Scripts" para procesar todos los archivos `.mq5` dentro de esa carpeta y sus subdirectorios.
-   **Para la carpeta actual:** Haz clic derecho en el fondo de una carpeta en el Explorador y selecciona "Fix MQ5 Scripts" para procesar todos los archivos `.mq5` en el directorio actual y sus subdirectorios.

### Modo de Línea de Comandos

Abre una terminal y ejecuta el archivo con la ruta al archivo o un patrón glob para múltiples archivos.

```sh
./fix-SQ-scripts.exe "ruta/a/tu/archivo.mq5"
./fix-SQ-scripts.exe "*.mq5"
```

## Desinstalación

Ejecuta el script `uninstall.ps1` para eliminar los accesos directos del menú contextual. Este script se genera automáticamente cuando ejecutas el instalador.

## Detalles Técnicos

La aplicación está escrita en Go y utiliza la biblioteca Fyne para la GUI. El script `install.ps1` crea entradas en el registro para agregar los accesos directos del menú contextual.

### Corrección del Menú Contextual

El script `install.ps1` original usaba `"%1"` para pasar la ruta a la aplicación para todas las entradas del menú contextual. Esto funciona para archivos y carpetas, pero no para el menú contextual del fondo de la carpeta. El script corregido ahora usa `"%V"` para el fondo de la carpeta, que pasa correctamente la ruta del directorio actual a la aplicación.