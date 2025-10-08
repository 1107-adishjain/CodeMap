// tools/main.js
const { analyzeDirectory } = require('./src/file_processor');

// Get the target directory from the command-line arguments.
// process.argv[0] is 'node', process.argv[1] is 'main.js'.
const targetDir = process.argv[2];

if (!targetDir) {
  console.error('Error: Please provide a directory to analyze.');
  process.exit(1);
}

// --- Main Execution Logic ---
try {
  // Call the single, powerful function from file_processor.js
  const analysisResults = analyzeDirectory(targetDir);

  const finalOutput = {
    files: analysisResults,
  };

  // The most important step: Print the final JSON to standard output.
  // The Go backend will capture this output.
  console.log(JSON.stringify(finalOutput, null, 2));

} catch (error) {
  console.error(`NODE.JS ANALYZER ERROR: ${error.message}`);
  console.error(`STACK TRACE: ${error.stack}`);
  process.exit(1);
}
