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
	
	let buffer;
	try {
		buffer = fs.readFileSync(target);
	} catch(readError) {
		throw new Error('Failed to read file: ' + String(readError));
	}
	
	let base64Content;
	try {
		base64Content = buffer.toString('base64');
	} catch(encodeError) {
		throw new Error('Failed to encode file to base64: ' + String(encodeError));
	}
	
	console.log(JSON.stringify({
		ok: true,
		path: target,
		base64: base64Content,
		size: buffer.length
	}));
} catch(error) {
	console.log(JSON.stringify({
		ok: false,
		error: String(error.message || error)
	}));
}

