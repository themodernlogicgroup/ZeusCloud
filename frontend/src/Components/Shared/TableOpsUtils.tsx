import { TableRow } from './Table';

export const searchFilterFn = (allRows: Array<TableRow>, filterValue: string): Array<TableRow> => {
    if(filterValue === "") {
        return allRows;
    }

    return allRows.filter(row => {
        for (var col of row.columns) {
            if (col.accessor_key === "description" && typeof(col.value) === "string"&&
                col.value.includes(filterValue)) {
                return true;
            }

            if (col.accessor_key === "service" && typeof(col.value) === "string" &&
                col.value.includes(filterValue)) {
                return true;
            }

            if (col.accessor_key === "risk_categories" && col.value instanceof Array && 
                col.value.some(v => v.includes(filterValue))) {
                return true;
            }

            if (col.accessor_key === "severity" && col.value === filterValue) {
                return true;
            }
        }
    });
}

export const risksFilterFn = (allRows: Array<TableRow>, filterValue: string): Array<TableRow> => {
    if (filterValue === "" || filterValue === "All") {
        return allRows;
    }

    return allRows.filter(row => {
        for (var col of row.columns) {
            if (col.accessor_key === "risk_categories" && col.value instanceof Array && 
                col.value.indexOf(filterValue) > -1) {
                return true;
            }
        }
    });
}

export const severityFilterFn = (allRows: Array<TableRow>, filterValue: string): Array<TableRow> => {
    if (filterValue === "" || filterValue === "All") {
        return allRows;
    }

    return allRows.filter(row => {
        for (var col of row.columns) {
            if (col.accessor_key === "severity" && col.value === filterValue) {
                return true;
            }
        }
    });
}

export const alertNumberSortTypeFn = (rowA: TableRow, rowB: TableRow): number => {
    var a: number = 0; 
    var b: number = 0; 
    for(var column of rowA.columns) {
        if(column.accessor_key === "alerts" && typeof(column.value) === "number") { 
           a = column.value
           break;
        }
    }
    
    for(var column of rowB.columns) {
        if(column.accessor_key === "alerts" && typeof(column.value) === "number") {
           b = column.value
           break;
        }
    } 
    console.log(a,b)
    if (a === b) {
        return 0;
    } else if (a > b) {
        return 1;
    } else {
        return -1;
    }
}

export const severitySortTypeFn = (rowA: TableRow, rowB: TableRow): number => {

    var a; 
    var b; 
    for(var column of rowA.columns) {
        if(column.accessor_key === "severity") { 
           a = column.value
           break;
        }
    }
    
    for(var column of rowB.columns) {
        if(column.accessor_key === "severity") {
           b = column.value
           break;
        }
    } 
    if (a === b) {
        return 0;
    }

    if (a === "Low") {
        return -1;
    }

    if (a === "Moderate") {
        if (b === "High" || b === "Critical") {
            return -1;
        }
        return 1;
    }

    if (a === "High") {
        if (b === "Critical") {
            return -1;
        }
        return 1;
    }
    return 1;
};