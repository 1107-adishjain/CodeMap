const Go = require('tree-sitter-go');
const Python = require('tree-sitter-python');
const Cpp = require('tree-sitter-cpp');
const Java = require('tree-sitter-java');
const { typescript, tsx } = require('tree-sitter-typescript');
const Html = require('tree-sitter-html');
const Css = require('tree-sitter-css');
const Kotlin = require('tree-sitter-kotlin');
const Swift = require('tree-sitter-swift');
const Dart = require('tree-sitter-dart');
const Json = require('tree-sitter-json');
const Yaml = require('tree-sitter-yaml');

const goExtractor = require('./extractors/go_extractor');
const pyExtractor = require('./extractors/py_extractor');
const jsExtractor = require('./extractors/js_extractor');
const cppExtractor = require('./extractors/cpp_extractor');
const javaExtractor = require('./extractors/java_extractor');
const tsExtractor = require('./extractors/ts_extractor');
const htmlExtractor = require('./extractors/html_extractor');
const cssExtractor = require('./extractors/css_extractor');
const kotlinExtractor = require('./extractors/kotlin_extractor');
const swiftExtractor = require('./extractors/swift_extractor');
const dartExtractor = require('./extractors/dart_extractor');
const jsonExtractor = require('./extractors/json_extractor');
const yamlExtractor = require('./extractors/yaml_extractor');

const ignoreExtensions = new Set([
    // Models and Assets
    '.glb', '.gltf', '.obj', '.fbx',
    // Images
    '.svg', '.ico', '.png', '.jpg', '.jpeg', '.gif', '.webp', '.bmp',
    // Fonts
    '.woff', '.woff2', '.ttf', '.otf', '.eot',
    // Documents
    '.md', '.txt', '.log', '.pdf', '.doc', '.docx', '.ppt', '.pptx',
    // Configuration & Lockfiles (that are not parsed)
    '.lock', '.env',
    // Git
    '.gitignore', '.gitattributes', '.gitmodules',
    // System / Temp
    '.DS_Store', '.bak', '.tmp', '.cache',
    // Archives
    '.zip', '.tar', '.gz', '.7z', '.rar',
    // Binaries
    '.exe', '.dll', '.bin', '.so', '.jar',
    // Media
    '.mp3', '.mp4', '.mov', '.avi', '.wav',
]);

// Directories to completely ignore during analysis
const ignoreDirs = new Set([
    '.git', 
    'node_modules', 
    'temp-uploads', 
    'temp-clones',
    '.next',           // Next.js build directory
    'dist',           // Common build output
    'build',          // Common build output  
    '.nuxt',          // Nuxt.js build directory
    'coverage',       // Test coverage
    '.nyc_output',    // NYC coverage
    'target',         // Rust/Java build output
    'bin',            // Binary output
    'obj',            // C# build output
    '.vscode',        // VS Code settings
    '.idea',          // IntelliJ settings
    '__pycache__',    // Python cache
    '.pytest_cache'   // Pytest cache
]);

const languageConfig = {
    '.js': { grammar: tsx, extractor: jsExtractor },
    '.mjs': { grammar: tsx, extractor: jsExtractor },
    '.jsx': { grammar: tsx, extractor: jsExtractor },
    '.ts': { grammar: typescript, extractor: tsExtractor },
    '.tsx': { grammar: tsx, extractor: tsExtractor },
    '.go': { grammar: Go, extractor: goExtractor },
    '.py': { grammar: Python, extractor: pyExtractor },
    '.cpp': { grammar: Cpp, extractor: cppExtractor },
    '.h': { grammar: Cpp, extractor: cppExtractor },
    '.hpp': { grammar: Cpp, extractor: cppExtractor },
    '.java': { grammar: Java, extractor: javaExtractor },
    '.html': { grammar: Html, extractor: htmlExtractor },
    '.css': { grammar: Css, extractor: cssExtractor },
    '.kt': { grammar: Kotlin, extractor: kotlinExtractor },
    '.kts': { grammar: Kotlin, extractor: kotlinExtractor },
    '.swift': { grammar: Swift, extractor: swiftExtractor },
    '.dart': { grammar: Dart, extractor: dartExtractor },
    '.json': { grammar: Json, extractor: jsonExtractor },
    '.yaml': { grammar: Yaml, extractor: yamlExtractor },
    '.yml': { grammar: Yaml, extractor: yamlExtractor },
};

module.exports = { languageConfig, ignoreExtensions, ignoreDirs };
