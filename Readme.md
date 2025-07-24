# Corrector de Scripts de SQ

Esta es una herramienta simple para parchear archivos `.mq5`, reemplazando una línea de código específica con una versión corregida. Se puede ejecutar desde la línea de comandos o a través de una interfaz gráfica de usuario (GUI) instalando accesos directos en el menú contextual.

## El Problema

Los scripts `.mq5` originales contienen un cálculo incorrecto para `PointValue`. Esta herramienta reemplaza la línea defectuosa con una versión corregida que utiliza `SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_CONTRACT_SIZE)`.

## Características

- **Interfaz de Línea de Comandos (CLI):** Procesa archivos directamente desde la terminal.
- **Interfaz Gráfica de Usuario (GUI):** Una ventana simple que muestra el estado de cada archivo que se está procesando.
- **Integración con Menú Contextual:** Haz clic derecho en archivos `.mq5`, carpetas o en el fondo de una carpeta en el Explorador de Windows para parchear los archivos.

## Requisitos Previos

- **Para todos los usuarios:** Se necesita **PowerShell** para ejecutar los scripts de instalación (`install.ps1`) y desinstalación (`uninstall.ps1`). PowerShell viene preinstalado en las versiones modernas de Windows.
- **Solo para desarrolladores:** Se necesita **Go** si deseas compilar la aplicación desde el código fuente.

## Instalación

Se proporcionan dos métodos de instalación: uno para usuarios finales que solo desean utilizar la herramienta y otro para desarrolladores que desean compilarla desde el código fuente.

### Para Usuarios Finales

No es necesario instalar Go. El ejecutable (`fix-SQ-scripts.exe`) ya está incluido.

1.  Asegúrate de tener PowerShell (generalmente incluido en Windows).
2.  Haz clic derecho en `install.ps1` y selecciona "Ejecutar con PowerShell" para agregar los accesos directos del menú contextual.

### Para Desarrolladores (Compilar desde fuente)

Si deseas modificar o compilar el programa, necesitarás Go.

1.  Ejecuta `build.ps1` para compilar el proyecto. Esto creará/reemplazará `fix-SQ-scripts.exe`.
2.  Ejecuta `install.ps1` para registrar los comandos del menú contextual para tu nueva versión.

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

La aplicación está escrita en Go y utiliza la biblioteca Fyne para la GUI. Una vez compilado, el ejecutable `fix-SQ-scripts.exe` es autónomo y no requiere dependencias externas para ejecutarse. El script `install.ps1` crea entradas en el registro para agregar los accesos directos del menú contextual.

### Corrección del Menú Contextual

El script `install.ps1` original usaba `"%1"` para pasar la ruta a la aplicación para todas las entradas del menú contextual. Esto funciona para archivos y carpetas, pero no para el menú contextual del fondo de la carpeta. El script corregido ahora usa `"%V"` para el fondo de la carpeta, que pasa correctamente la ruta del directorio actual a la aplicación.

## Fix for `partialScript.mq5`

A specific issue was identified in the `partialScript.mq5` script where the `sqMMFixedAmount` function was using an incorrect calculation for determining the lot size. This led to order placement failures with "Unknown error" (4307).

The original calculation used a naive point-difference approach which did not account for the instrument's contract specifications. The fix involved modifying the `sqMMFixedAmount` function to use `OrderCalcProfit` for accurate drawdown calculation. This ensures that the lot size is calculated correctly according to the broker's specifications.

The patched script `partialScript-148076.mq5` contains this fix and has been verified to resolve the order placement issue.