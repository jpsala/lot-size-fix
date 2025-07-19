# Problema: Error "Not Enough Money" (Dinero Insuficiente) en EAs generados por StrategyQuant para MT5

## 1. El Problema

Al utilizar Asesores Expertos (EAs) generados por StrategyQuant para operar CFDs de índices (por ejemplo, NDX100) en MetaTrader 5, con frecuencia ocurre un error de "Not enough money" (Dinero insuficiente), impidiendo que se abran operaciones. Este problema puede aparecer en un bróker y en otro no, incluso con tamaños de cuenta y configuraciones de riesgo idénticas.

La causa raíz es un cálculo incorrecto del tamaño del lote dentro del código MQL5 generado. El EA calcula un tamaño de lote que es demasiado grande para el margen disponible de la cuenta porque malinterpreta las propiedades del instrumento en ciertos brókers.

## 2. La Causa Raíz: Configuraciones Inconsistentes del Bróker

Este problema se origina en un cálculo frágil en el módulo de gestión de capital de StrategyQuant que no tiene en cuenta las inconsistencias en cómo los diferentes brókers configuran sus instrumentos de CFD.

El código generado calcula el valor del punto del instrumento (`PointValue`) usando la siguiente fórmula:

```mql5
double PointValue = SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_VALUE) / SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_SIZE);
```

Esta fórmula no es confiable. Asume que el `SYMBOL_TRADE_TICK_VALUE` de un bróker siempre estará perfectamente escalado con el `SYMBOL_TRADE_CONTRACT_SIZE` del instrumento. Esta no es una suposición segura.

### Ejemplo del Mundo Real: Darwinex vs. FundedNext

Una comparación del instrumento NDX100 en dos brókers diferentes demuestra el problema:

| Propiedad | Darwinex | FundedNext |
| :--- | :--- | :--- |
| **`Contract size`** | **10** | **10** |
| `Tick size` | 0.1 | 0.01 |
| `Tick value` | 1.0 (inferido) | 0.01 |

- **En Darwinex**, la fórmula funciona por coincidencia: `PointValue = 1.0 / 0.1 = 10`. Esto coincide con el `Contract size`, por lo que el cálculo del lote es correcto.
- **En FundedNext**, la fórmula falla: `PointValue = 0.01 / 0.01 = 1`. El script calcula un `PointValue` de $1, cuando el valor real basado en el `Contract size` es de $10.

Como resultado, en FundedNext, el EA subestima el riesgo por un factor de 10 y calcula un tamaño de lote que es **10 veces demasiado grande**, causando el error "Not enough money".

## 3. La Solución

### La Corrección del Código

La solución es modificar el cálculo de `PointValue` para usar la única propiedad confiable para el cálculo de riesgo: el tamaño del contrato del símbolo.

Cambie esta línea en el archivo `.mq5`:

**De:**
```mql5
double PointValue = SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_VALUE) / SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_TICK_SIZE);
```

**A:**
```mql5
double PointValue = SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_CONTRACT_SIZE);
```

Esto asegura que el riesgo se calcule correctamente en cualquier bróker, independientemente de la configuración de su servidor.

## 4. Explicación de Términos Técnicos

Para entender completamente el problema y la solución, es crucial conocer qué representa cada uno de estos términos en MQL5.

### `SYMBOL_TRADE_CONTRACT_SIZE`
- **Qué es:** Es el tamaño del contrato de un instrumento, expresado en la moneda base del activo. Representa la cantidad de unidades del activo subyacente que se negocian en un lote estándar (1.0). Este es el valor más importante y fiable para determinar el valor real de una operación.
- **Ejemplo:** Si el `SYMBOL_TRADE_CONTRACT_SIZE` para el índice NDX100 es **10**, significa que un lote estándar (1.0) de NDX100 controla un valor nominal de 10 veces el precio del índice. Si el NDX100 cotiza a 18,000, el valor total de un lote es `10 * 18,000 = $180,000`.

### `SYMBOL_TRADE_TICK_SIZE`
- **Qué es:** Es el cambio mínimo de precio posible para un símbolo. Se le conoce comúnmente como "tick".
- **Ejemplo:** Si el `SYMBOL_TRADE_TICK_SIZE` del NDX100 es **0.01**, el precio del índice solo puede moverse en incrementos de 0.01 (por ejemplo, de 18,000.00 a 18,000.01).

### `SYMBOL_TRADE_TICK_VALUE`
- **Qué es:** Es el valor monetario de un solo "tick" de movimiento para un lote estándar (1.0). Este valor es calculado y proporcionado por el bróker.
- **Problema:** La inconsistencia de este valor entre brókers es la raíz del problema. Algunos brókers lo configuran en relación directa con el tamaño del contrato, mientras que otros no.
- **Ejemplo (FundedNext):** Con un `TICK_SIZE` de 0.01 y un `TICK_VALUE` de 0.01, el bróker está diciendo que un movimiento de 0.01 en el precio resulta en una ganancia/pérdida de $0.01 por lote. Esto es incorrecto para un contrato de tamaño 10.
- **Ejemplo (Darwinex):** Con un `TICK_SIZE` de 0.1 y un `TICK_VALUE` de 1.0, el cálculo es `1.0 / 0.1 = 10`, que coincide con el tamaño del contrato. Aquí funciona, pero por casualidad.

### `PointValue` (Valor del Punto)
- **Qué es:** En el contexto del código de StrategyQuant, `PointValue` es una variable interna utilizada para calcular el tamaño del lote. El objetivo de esta variable es determinar cuánto vale un movimiento de un punto completo en el precio del instrumento.
- **Cálculo Erróneo:** `PointValue = SYMBOL_TRADE_TICK_VALUE / SYMBOL_TRADE_TICK_SIZE`. Como se demostró, esta fórmula es propensa a errores porque depende del `TICK_VALUE` configurado por el bróker.
- **Cálculo Correcto:** Al asignar directamente `PointValue = SYMBOL_TRADE_CONTRACT_SIZE`, eliminamos la dependencia de las configuraciones `TICK_VALUE` y `TICK_SIZE`. Usamos el valor más robusto y estandarizado (`CONTRACT_SIZE`) para asegurar que el cálculo del riesgo y del lotaje sea siempre correcto, sin importar el bróker.

## 5. Herramienta de Parcheo Automatizado

Para simplificar la aplicación de esta corrección, se proporciona la herramienta `fix-SQ-scripts.exe`. Esta herramienta busca y reemplaza automáticamente la línea incorrecta en cualquier archivo `.mq5` dentro del directorio actual y sus subdirectorios.

**Cómo usar la herramienta:**

1.  Abra una terminal.
2.  Ejecute el programa con un patrón de archivo. Por ejemplo, para parchear todos los archivos `.mq5` en el directorio actual:

    ```shell
    .\fix-SQ-scripts.exe "*.mq5"
    o
    .\fix-SQ-scripts.exe "c:\scripts\Strategy*.mq5"
    ```

Esto parcheará los archivos, dejándolos listos para ser compilados en MetaEditor con el cálculo de tamaño de lote correcto.

### El repositorio está aquí: https://github.com/jpsala/lot-size-fix