const Parser = require('tree-sitter');
const fs = require('fs');
const path = require('path');
const { languageConfig, ignoreExtensions, ignoreDirs } = require('./language');

const parser = new Parser();

function analyzeFile(filePath) {
    const extension = path.extname(filePath);

    // Skip files based on extension or if they have no extension
    if (ignoreExtensions.has(extension) || extension === '') {
        return null;
    }

    const config = languageConfig[extension];

    // Skip files for which we have no defined language configuration
    if (!config) {
        // Silently skip unsupported files to keep output clean
        return null;
    }

    try {
        const fileContent = fs.readFileSync(filePath, 'utf8');
        parser.setLanguage(config.grammar);
        const tree = parser.parse(fileContent);
        const analysis = config.extractor.extract(tree, config);

        return {
            path: filePath,
            language: config.grammar.name,
            ...analysis,
        };
    } catch (err) {
        // Silently skip files that can't be parsed to keep JSON output clean
        // Don't return error objects as they clutter the results
        return null;
    }
}

function analyzeDirectory(directoryPath) {
    let results = [];
    let items;

    try {
        items = fs.readdirSync(directoryPath);
    } catch (err) {
        // Silently skip inaccessible directories
        return [];
    }

    for (const item of items) {
        // Check if the directory itself should be ignored
        if (ignoreDirs.has(item)) {
            continue; // Skip this directory entirely
        }

        const itemPath = path.join(directoryPath, item);
        let stat;

        try {
            stat = fs.statSync(itemPath);
        } catch (err) {
            // Silently skip inaccessible files
            continue;
        }

        if (stat.isDirectory()) {
            results = results.concat(analyzeDirectory(itemPath));
        } else {
            const fileAnalysis = analyzeFile(itemPath);
            if (fileAnalysis) {
                results.push(fileAnalysis);
            }
        }
    }
    return results;
}

module.exports = { analyzeDirectory };
