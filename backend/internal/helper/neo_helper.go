package helper

import (
	"codemap/backend/internal/models"
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Helper functions for the transaction
func CreateNodesForFile(ctx context.Context, tx neo4j.ManagedTransaction, file models.File) error {
	// Create File node
	_, err := tx.Run(ctx, `
        MERGE (f:File {path: $path})
        ON CREATE SET f.language = $language
    `, map[string]any{
		"path":     file.Path,
		"language": file.Language,
	})
	if err != nil {
		return err
	}

	// Create Class and Property nodes
	for _, class := range file.Classes {
		classID := fmt.Sprintf("%s#%s", file.Path, class.Name)
		_, err := tx.Run(ctx, `
            MATCH (f:File {path: $filePath})
            MERGE (c:Class {id: $classID})
            ON CREATE SET c.name = $name, c.is_exported = $is_exported
            MERGE (f)-[:CONTAINS]->(c)
        `, map[string]any{
			"filePath":    file.Path,
			"classID":     classID,
			"name":        class.Name,
			"is_exported": class.IsExported,
		})
		if err != nil {
			return err
		}

		for _, propName := range class.Properties {
			propID := fmt.Sprintf("%s::%s", classID, propName)
			_, err := tx.Run(ctx, `
                MATCH (c:Class {id: $classID})
                MERGE (p:Property {id: $propID})
                ON CREATE SET p.name = $name
                MERGE (c)-[:HAS_PROPERTY]->(p)
            `, map[string]any{
				"classID": classID,
				"propID":  propID,
				"name":    propName,
			})
			if err != nil {
				return err
			}
		}
	}

	// Create Function and Parameter nodes with enhanced properties
	for _, function := range file.Functions {
		funcID := fmt.Sprintf("%s#%s", file.Path, function.Name)
		_, err := tx.Run(ctx, `
            MATCH (f:File {path: $filePath})
            MERGE (fn:Function {id: $funcID})
            ON CREATE SET 
                fn.name = $name, 
                fn.is_exported = $is_exported, 
                fn.is_method_of = $is_method_of,
                fn.return_types = $return_types,
                fn.param_count = $param_count
            MERGE (f)-[:CONTAINS]->(fn)
        `, map[string]any{
			"filePath":     file.Path,
			"funcID":       funcID,
			"name":         function.Name,
			"is_exported":  function.IsExported,
			"is_method_of": function.IsMethodOf,
			"return_types": function.ReturnTypes,
			"param_count":  len(function.Params),
		})
		if err != nil {
			return err
		}

		// Create Parameter nodes with enhanced information
		for i, paramName := range function.Params {
			paramID := fmt.Sprintf("%s(%s)", funcID, paramName)
			_, err := tx.Run(ctx, `
                MATCH (fn:Function {id: $funcID})
                MERGE (p:Parameter {id: $paramID})
                ON CREATE SET 
                    p.name = $name,
                    p.position = $position
                MERGE (fn)-[:HAS_PARAMETER]->(p)
            `, map[string]any{
				"funcID":   funcID,
				"paramID":  paramID,
				"name":     paramName,
				"position": i + 1,
			})
			if err != nil {
				return err
			}
		}
	}

	// Create Import nodes and relationships between files
	for _, importItem := range file.Imports {
		importID := fmt.Sprintf("%s->%s", file.Path, importItem.Source)
		_, err := tx.Run(ctx, `
            MATCH (f:File {path: $filePath})
            MERGE (imp:Import {id: $importID})
            ON CREATE SET 
                imp.source = $source,
                imp.from_file = $filePath
            MERGE (f)-[:HAS_IMPORT]->(imp)
        `, map[string]any{
			"filePath": file.Path,
			"importID": importID,
			"source":   importItem.Source,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateRelationshipsForFile(ctx context.Context, tx neo4j.ManagedTransaction, file models.File) error {
	// Create enhanced IMPORTS relationships between files
	for _, imp := range file.Imports {
		if imp.Source != "" {
			_, err := tx.Run(ctx, `
                MATCH (importer:File {path: $importerPath})
                // Try to find exact file matches first, then partial matches
                OPTIONAL MATCH (imported:File) WHERE imported.path ENDS WITH $importSource OR imported.path CONTAINS $importSource
                WITH importer, imported, $importSource as source
                FOREACH (f IN CASE WHEN imported IS NOT NULL THEN [imported] ELSE [] END |
                    MERGE (importer)-[:IMPORTS {
                        source: source, 
                        import_type: 'internal',
                        resolved: true
                    }]->(f)
                )
                // Create external dependency node for unresolved imports (libraries, etc.)
                FOREACH (x IN CASE WHEN imported IS NULL THEN [1] ELSE [] END |
                    MERGE (ext:ExternalDependency {name: source, type: 'library'})
                    MERGE (importer)-[:DEPENDS_ON {
                        source: source, 
                        import_type: 'external',
                        resolved: false
                    }]->(ext)
                )
            `, map[string]any{
				"importerPath": file.Path,
				"importSource": imp.Source,
			})
			if err != nil {
				fmt.Printf("Warning: Could not create import relationship: %v\n", err)
			}
		}
	}

	// Create enhanced class-method relationships
	for _, class := range file.Classes {
		classID := fmt.Sprintf("%s#%s", file.Path, class.Name)

		// Create direct relationships to class methods from class.Methods array
		for _, methodName := range class.Methods {
			methodID := fmt.Sprintf("%s#%s", file.Path, methodName)
			_, err := tx.Run(ctx, `
                MATCH (c:Class {id: $classID})
                MATCH (fn:Function {id: $methodID})
                MERGE (c)-[:HAS_METHOD {method_name: $methodName}]->(fn)
            `, map[string]any{
				"classID":    classID,
				"methodID":   methodID,
				"methodName": methodName,
			})
			if err != nil {
				fmt.Printf("Warning: Could not create class-method relationship: %v\n", err)
			}
		}
	}

	// Create enhanced function relationships
	for _, function := range file.Functions {
		funcID := fmt.Sprintf("%s#%s", file.Path, function.Name)

		// Method ownership relationships
		if function.IsMethodOf != "" {
			classID := fmt.Sprintf("%s#%s", file.Path, function.IsMethodOf)
			_, err := tx.Run(ctx, `
                MATCH (c:Class {id: $classID})
                MATCH (fn:Function {id: $funcID})
                MERGE (c)-[:OWNS_METHOD]->(fn)
                SET fn.is_method = true
            `, map[string]any{
				"classID": classID,
				"funcID":  funcID,
			})
			if err != nil {
				fmt.Printf("Warning: Could not create method ownership: %v\n", err)
			}
		}

		// Enhanced function call relationships
		for i, calledFuncName := range function.Calls {
			_, err := tx.Run(ctx, `
                MATCH (caller:Function {id: $callerID})
                // Try to find called function in same file first, then globally
                OPTIONAL MATCH (callee_same_file:Function) 
                WHERE callee_same_file.id STARTS WITH $sameFilePrefix AND callee_same_file.name = $calleeName
                OPTIONAL MATCH (callee_global:Function {name: $calleeName})
                WHERE callee_same_file IS NULL
                WITH caller, COALESCE(callee_same_file, callee_global) as callee, $calleeName as funcName
                FOREACH (f IN CASE WHEN callee IS NOT NULL THEN [callee] ELSE [] END |
                    MERGE (caller)-[:CALLS {
                        function_name: funcName, 
                        call_order: $callOrder,
                        call_type: CASE WHEN callee.is_method_of IS NOT NULL THEN 'method' ELSE 'function' END
                    }]->(f)
                )
            `, map[string]any{
				"callerID":       funcID,
				"calleeName":     calledFuncName,
				"callOrder":      i + 1,
				"sameFilePrefix": file.Path + "#",
			})
			if err != nil {
				fmt.Printf("Warning: Could not create CALLS relationship from %s to %s: %v\n", funcID, calledFuncName, err)
			}
		}

		// Create return type relationships for type analysis
		for _, returnType := range function.ReturnTypes {
			if returnType != "" && returnType != "void" && returnType != "any" {
				_, err := tx.Run(ctx, `
                    MATCH (fn:Function {id: $funcID})
                    MERGE (rt:ReturnType {name: $returnType})
                    MERGE (fn)-[:RETURNS]->(rt)
                `, map[string]any{
					"funcID":     funcID,
					"returnType": returnType,
				})
				if err != nil {
					fmt.Printf("Warning: Could not create return type relationship: %v\n", err)
				}
			}
		}
	}
	return nil
}
