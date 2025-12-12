(function(){
try{
const cmd={{CMD}};
const exec=require('child_process').exec;
const timeout=30000;
let output='';
let success=true;
let completed=false;
const childProcess=exec(cmd,{maxBuffer:104857600,encoding:'utf8'},function(err,stdout,stderr){
if(completed)return;
completed=true;
if(err){
success=false;
output=err.message||err.toString();
if(stdout)output+='\n'+stdout;
if(stderr)output+='\n'+stderr;
}else{
output=stdout||'';
if(stderr)output+=stderr;
}
});
const timeoutId=setTimeout(function(){
if(!completed){
completed=true;
try{childProcess.kill();}catch(e){}
success=false;
output='Command execution timeout after '+(timeout/1000)+' seconds';
}
},timeout);
const startTime=Date.now();
const maxWaitTime=timeout+5000;
while(!completed&&(Date.now()-startTime)<maxWaitTime){
const waitStart=Date.now();
while(Date.now()-waitStart<10){}
}
clearTimeout(timeoutId);
if(!completed){
try{childProcess.kill();}catch(e){}
success=false;
output='Command execution timeout';
}
const jsonStr=JSON.stringify({output:output,success:success});
const outputBase64=Buffer.from(jsonStr,'utf8').toString('base64');
return JSON.stringify({success:true,output_base64:outputBase64});
}catch(err){
const errorMsg=err.message||err.toString();
return JSON.stringify({success:false,output_base64:'',error:errorMsg});
}
})();
