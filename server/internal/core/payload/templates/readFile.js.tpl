const fs = require('node:fs');
const path = require('node:path');
try {
	let filePath = {{FILE_PATH}};
	if (!filePath || typeof filePath !== 'string') {
		throw new Error('Invalid file path');
	}
	
	let target;
	try {
		target = path.resolve(filePath);
	} catch(pathError) {
		throw new Error('Failed to resolve file path: ' + String(pathError));
	}
	
	if (!fs.existsSync(target)) {
		throw new Error('File not found: ' + target);
	}
	
	let stats;
	try {
		stats = fs.statSync(target);
	} catch(statError) {
		throw new Error('Failed to get file stats: ' + String(statError));
	}
	
	if (!stats.isFile()) {
		throw new Error('Path is not a file: ' + target);
	}
	
	let content;
	try {
		content = fs.readFileSync(target, 'utf8');
	} catch(readError) {
		throw new Error('Failed to read file: ' + String(readError));
	}
	
	console.log(JSON.stringify({
		ok: true,
		path: target,
		content: content,
		size: content.length
	}));
} catch(error) {
	console.log(JSON.stringify({
		ok: false,
		error: String(error.message || error)
	}));
}

