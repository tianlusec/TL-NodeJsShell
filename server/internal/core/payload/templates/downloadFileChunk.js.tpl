const fs = require('node:fs');
const path = require('node:path');
try {
	let filePath = {{FILE_PATH}};
	let chunkIndex = {{CHUNK_INDEX}};
	let chunkSize = {{CHUNK_SIZE}};
	
	if (!filePath || typeof filePath !== 'string') {
		throw new Error('Invalid file path');
	}
	
	if (typeof chunkIndex !== 'number' || chunkIndex < 0) {
		throw new Error('Invalid chunk index');
	}
	
	if (typeof chunkSize !== 'number' || chunkSize <= 0) {
		throw new Error('Invalid chunk size');
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
	
	let fileSize = stats.size;
	let startPos = chunkIndex * chunkSize;
	let endPos = Math.min(startPos + chunkSize, fileSize);
	
	if (startPos >= fileSize) {
		throw new Error('Chunk index out of range');
	}
	
	// 读取指定范围的数据
	let fd;
	try {
		fd = fs.openSync(target, 'r');
	} catch(openError) {
		throw new Error('Failed to open file: ' + String(openError));
	}
	
	let buffer;
	try {
		buffer = Buffer.alloc(endPos - startPos);
		let bytesRead = fs.readSync(fd, buffer, 0, endPos - startPos, startPos);
		if (bytesRead !== endPos - startPos) {
			buffer = buffer.slice(0, bytesRead);
		}
	} catch(readError) {
		fs.closeSync(fd);
		throw new Error('Failed to read file chunk: ' + String(readError));
	}
	
	fs.closeSync(fd);
	
	// Base64编码
	let base64Content;
	try {
		base64Content = buffer.toString('base64');
	} catch(encodeError) {
		throw new Error('Failed to encode chunk to base64: ' + String(encodeError));
	}
	
	// 计算总 chunk 数
	let totalChunks = Math.ceil(fileSize / chunkSize);
	
	// 同时返回下划线命名，便于前端使用
	console.log(JSON.stringify({
		ok: true,
		path: target,
		// camelCase
		chunkIndex: chunkIndex,
		totalChunks: totalChunks,
		chunkSize: buffer.length,
		fileSize: fileSize,
		base64: base64Content,
		// snake_case（前端按此字段读取）
		chunk_index: chunkIndex,
		total_chunks: totalChunks,
		chunk_size: buffer.length,
		file_size: fileSize,
		data: base64Content
	}));
} catch(error) {
	console.log(JSON.stringify({
		ok: false,
		error: String(error.message || error),
		chunkIndex: typeof chunkIndex !== 'undefined' ? chunkIndex : -1
	}));
}


