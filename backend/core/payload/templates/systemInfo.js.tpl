const os = require('node:os');
const fs = require('node:fs');
const path = require('node:path');
try {
	let platform, arch, hostname, type, release, uptime, totalmem, freemem, cpus;
	
	try {
		platform = os.platform();
	} catch(e) { platform = 'unknown'; }
	
	try {
		arch = os.arch();
	} catch(e) { arch = 'unknown'; }
	
	try {
		hostname = os.hostname();
	} catch(e) { hostname = 'unknown'; }
	
	try {
		type = os.type();
	} catch(e) { type = 'unknown'; }
	
	try {
		release = os.release();
	} catch(e) { release = 'unknown'; }
	
	try {
		uptime = os.uptime();
	} catch(e) { uptime = 0; }
	
	try {
		totalmem = os.totalmem();
	} catch(e) { totalmem = 0; }
	
	try {
		freemem = os.freemem();
	} catch(e) { freemem = 0; }
	
	try {
		let cpuList = os.cpus();
		cpus = cpuList ? cpuList.length : 0;
	} catch(e) { cpus = 0; }
	
	let userInfo = null;
	try {
		userInfo = os.userInfo();
	} catch(userError) {
		userInfo = {
			uid: -1,
			gid: -1,
			username: 'unknown',
			homedir: process.env.HOME || process.env.USERPROFILE || '/',
			shell: process.env.SHELL || process.env.COMSPEC || null
		};
	}
	
	let envVars = {};
	try {
		envVars = process.env || {};
	} catch(e) {
		envVars = {};
	}
	
	let hostsContent = '';
	try {
		let hostsPath;
		if (platform === 'win32') {
			hostsPath = path.join(process.env.SYSTEMROOT || process.env.WINDIR || 'C:\\Windows', 'System32', 'drivers', 'etc', 'hosts');
		} else {
			hostsPath = '/etc/hosts';
		}
		if (fs.existsSync(hostsPath)) {
			try {
				hostsContent = fs.readFileSync(hostsPath, 'utf8');
			} catch(readError) {
				hostsContent = 'Failed to read hosts file: ' + String(readError);
			}
		} else {
			hostsContent = 'Hosts file not found: ' + hostsPath;
		}
	} catch(hostsError) {
		hostsContent = 'Failed to access hosts file: ' + String(hostsError);
	}
	
	let networkInterfaces = [];
	try {
		let interfaces = os.networkInterfaces();
		if (interfaces) {
			for (let name in interfaces) {
				if (interfaces.hasOwnProperty(name)) {
					let addrs = interfaces[name];
					if (addrs && Array.isArray(addrs)) {
						for (let i = 0; i < addrs.length; i++) {
							let addr = addrs[i];
							if (addr && addr.address) {
								let ipInfo = {
									interface: name,
									address: addr.address,
									family: addr.family || 'unknown',
									internal: addr.internal || false
								};
								if (addr.netmask) {
									ipInfo.netmask = addr.netmask;
								}
								if (addr.mac) {
									ipInfo.mac = addr.mac;
								}
								networkInterfaces.push(ipInfo);
							}
						}
					}
				}
			}
		}
	} catch(netError) {
		networkInterfaces = [];
	}
	
	const result = {
		ok: true,
		platform: platform,
		arch: arch,
		hostname: hostname,
		type: type,
		release: release,
		uptime: uptime,
		totalmem: totalmem,
		freemem: freemem,
		cpus: cpus,
		userInfo: userInfo,
		envVars: envVars,
		hosts: hostsContent,
		networkInterfaces: networkInterfaces
	};
	
	console.log(JSON.stringify(result));
} catch(error) {
	console.log(JSON.stringify({
		ok: false,
		error: String(error.message || error)
	}));
}

