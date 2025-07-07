import { QueryValue } from './types';

/**
 * Clase para manejar las filas de resultados de una query SQL
 * Equivalente a la estructura Rows en Go
 */
export class Rows {
  private columns: string[];
  private rows: any[][];
  private currentIndex: number = 0;

  constructor(columns: string[], rows: any[][]) {
    this.columns = columns;
    this.rows = rows;
  }

  /**
   * Retorna los nombres de las columnas
   */
  getColumns(): string[] {
    return [...this.columns]; // copia para evitar mutación
  }

  /**
   * Retorna todas las filas como array de objetos
   */
  getRows(): Record<string, QueryValue>[] {
    return this.rows.map(row => {
      const obj: Record<string, QueryValue> = {};
      for (let i = 0; i < this.columns.length; i++) {
        obj[this.columns[i]] = this.convertValue(row[i]);
      }
      return obj;
    });
  }

  /**
   * Retorna todas las filas como array de arrays
   */
  getRawRows(): any[][] {
    return this.rows.map(row => row.map(val => this.convertValue(val)));
  }

  /**
   * Indica si hay más filas disponibles
   */
  hasNext(): boolean {
    return this.currentIndex < this.rows.length;
  }

  /**
   * Avanza al siguiente registro y lo retorna como objeto
   */
  next(): Record<string, QueryValue> | null {
    if (!this.hasNext()) {
      return null;
    }

    const row = this.rows[this.currentIndex];
    const obj: Record<string, QueryValue> = {};
    
    for (let i = 0; i < this.columns.length; i++) {
      obj[this.columns[i]] = this.convertValue(row[i]);
    }

    this.currentIndex++;
    return obj;
  }

  /**
   * Reinicia el iterador al principio
   */
  reset(): void {
    this.currentIndex = 0;
  }

  /**
   * Retorna el número total de filas
   */
  length(): number {
    return this.rows.length;
  }

  /**
   * Convierte valores del servidor a tipos apropiados
   * Función auxiliar equivalente a convertValue en Go
   */
  private convertValue(val: any): QueryValue {
    if (val === null || val === undefined) {
      return null;
    }

    switch (typeof val) {
      case 'string':
        // Intentar convertir strings que representan números
        const intVal = parseInt(val, 10);
        if (!isNaN(intVal) && intVal.toString() === val) {
          return intVal;
        }
        
        const floatVal = parseFloat(val);
        if (!isNaN(floatVal) && floatVal.toString() === val) {
          return floatVal;
        }
        
        return val;
      
      case 'number':
        // JSON unmarshaling siempre devuelve number para números
        // Si es un entero, mantenerlo como entero
        if (Number.isInteger(val)) {
          return val;
        }
        return val;
      
      case 'boolean':
        return val;
      
      default:
        // Para otros tipos, convertir a string
        return String(val);
    }
  }
} 