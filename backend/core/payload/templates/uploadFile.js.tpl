const fs = require('node:fs');
const path = require('node:path');
try {
	let filePath = {{FILE_PATH}};
	if (!filePath || typeof filePath !== 'string') {
		throw new Error('Invalid file path');
	}
	
	let fileContent = {{FILE_CONTENT}};
	if (!fileContent || typeof fileContent !== 'string') {
		throw new Error('Invalid file content');
	}
	
	let target;
	try {
		target = path.resolve(filePath);
	} catch(pathError) {
		throw new Error('Failed to resolve file path: ' + String(pathError));
	}
	
	let content;
	try {
		content = Buffer.from(fileContent, 'base64');
	} catch(decodeError) {
		throw new Error('Failed to decode base64 content: ' + String(decodeError));
	}
	
	let dir;
	try {
		dir = path.dirname(target);
	} catch(dirError) {
		throw new Error('Failed to get directory path: ' + String(dirError));
	}
	
	try {
		if (!fs.existsSync(dir)) {
			try {
				fs.mkdirSync(dir, { recursive: true });
			} catch(mkdirError) {
				throw new Error('Failed to create directory: ' + String(mkdirError));
			}
		}
	} catch(dirCheckError) {
		throw new Error('Failed to check directory: ' + String(dirCheckError));
	}
	
	try {
		fs.writeFileSync(target, content);
	} catch(writeError) {
		throw new Error('Failed to write file: ' + String(writeError));
	}
	
	console.log(JSON.stringify({
		ok: true,
		path: target,
		size: content.length
	}));
} catch(error) {
	console.log(JSON.stringify({
		ok: false,
		error: String(error.message || error)
	}));
}

