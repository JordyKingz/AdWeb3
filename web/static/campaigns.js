'use strict';
const el = id => document.getElementById(id);
async function request(path, body) {
  try {
    const response = await fetch(path, {method: body ? 'POST' : 'GET', headers: {'Authorization':'Bearer '+el('credential').value,'Content-Type':'application/json','Idempotency-Key':crypto.randomUUID()},body:body?JSON.stringify(body):undefined});
    const text=await response.text();el('result').textContent=text;
    if(response.ok && body && path==='/v1/orders')el('order').value=JSON.parse(text).id;
  } catch(e) { el('result').textContent='Request failed. Check the connection and retry.'; }
}
el('create').onclick=()=>{const start=Math.floor(Date.now()/1000);request('/v1/orders',{payer:el('payer').value,token:el('token').value,chainId:Number(el('chain').value),amount:'25000000',start,end:start+604800});};
el('load').onclick=()=>request('/v1/orders/'+encodeURIComponent(el('order').value));
el('refundButton').onclick=()=>request('/v1/orders/'+encodeURIComponent(el('order').value)+'/refunds',{amount:el('refund').value});
