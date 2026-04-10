package servidor

import (
	"fmt"
	"strings"

	"github.com/flavio/flang/compiler/ast"
)

func (s *Servidor) renderHTML() string {
	theme := s.Program.Theme
	if theme == nil {
		theme = ast.DefaultTheme()
	}

	var b strings.Builder
	b.WriteString(`<!DOCTYPE html><html lang="pt-BR"><head><meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>` + cap(s.Program.System.Name) + `</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
<style>`)
	b.WriteString(s.cssV2(theme))
	b.WriteString(`</style></head><body class="dark">`)

	// Sidebar
	b.WriteString(`<aside class="sidebar" id="sidebar">`)
	b.WriteString(`<div class="sb-top">`)
	b.WriteString(`<div class="sb-brand"><div class="sb-logo">` + svgIcon("zap") + `</div>`)
	b.WriteString(`<span class="sb-name">` + cap(s.Program.System.Name) + `</span>`)
	b.WriteString(`<button class="sb-collapse" onclick="toggleCollapse()" title="Recolher">` + svgIcon("chevleft") + `</button></div>`)
	b.WriteString(`<nav class="sb-nav">`)
	b.WriteString(`<a class="sb-link active" onclick="irPara('dashboard',this)" href="#">`)
	b.WriteString(`<div class="sb-icon">` + svgIcon("layout") + `</div><span>Dashboard</span></a>`)
	for _, model := range s.Program.Models {
		name := lo(model.Name)
		icon := modelIcon(name)
		b.WriteString(fmt.Sprintf(`<a class="sb-link" onclick="irPara('%s',this)" href="#">`, name))
		b.WriteString(`<div class="sb-icon">` + svgIcon(icon) + `</div><span>` + cap(model.Name) + `</span></a>`)
	}
	b.WriteString(`</nav></div>`)
	// Sidebar footer
	b.WriteString(`<div class="sb-foot">`)
	b.WriteString(`<button class="sb-theme" onclick="toggleDark()">` + svgIcon("moon") + `<span>Tema</span></button>`)
	b.WriteString(`<div class="sb-powered">Flang v0.3</div>`)
	b.WriteString(`</div></aside>`)

	// Main
	b.WriteString(`<div class="main" id="main">`)
	// Topbar
	b.WriteString(`<header class="topbar glass">`)
	b.WriteString(`<button class="tb-menu" onclick="toggleSidebar()">` + svgIcon("menu") + `</button>`)
	b.WriteString(`<h1 id="page-title">Dashboard</h1>`)
	b.WriteString(`<div class="tb-end">`)
	b.WriteString(`<div class="tb-search glass-input"><input type="text" placeholder="Buscar..." id="global-search" oninput="buscaGlobal(this.value)">` + svgIcon("search") + `</div>`)
	b.WriteString(`</div></header>`)

	// Content
	b.WriteString(`<div class="content">`)

	// Dashboard
	b.WriteString(`<div id="secao-dashboard" class="section anim-in">`)
	// Bento grid stats
	b.WriteString(`<div class="bento">`)
	colors := []string{"#6366f1", "#8b5cf6", "#ec4899", "#f59e0b", "#10b981", "#3b82f6", "#ef4444", "#06b6d4"}
	for i, model := range s.Program.Models {
		name := lo(model.Name)
		icon := modelIcon(name)
		color := colors[i%len(colors)]
		b.WriteString(fmt.Sprintf(`<div class="bento-card" onclick="irParaNav('%s')" style="--accent:%s">`, name, color))
		b.WriteString(`<div class="bc-icon">` + svgIcon(icon) + `</div>`)
		b.WriteString(fmt.Sprintf(`<div class="bc-num" id="stat-%s">0</div>`, name))
		b.WriteString(`<div class="bc-label">` + cap(model.Name) + `</div>`)
		b.WriteString(`<div class="bc-glow"></div></div>`)
	}
	b.WriteString(`</div>`)

	// Activity feed
	b.WriteString(`<div class="dash-grid">`)
	b.WriteString(`<div class="card glass-card"><div class="card-head">` + svgIcon("activity") + `<h3>Atividade Recente</h3></div>`)
	b.WriteString(`<div id="atividade" class="ativ-list"><div class="empty-state">` + svgIcon("inbox") + `<p>Nenhuma atividade</p></div></div></div>`)
	// Quick info
	b.WriteString(`<div class="card glass-card"><div class="card-head">` + svgIcon("info") + `<h3>Informações</h3></div>`)
	b.WriteString(`<div class="info-list">`)
	b.WriteString(fmt.Sprintf(`<div class="info-row"><span class="info-k">Sistema</span><span class="info-v">%s</span></div>`, cap(s.Program.System.Name)))
	b.WriteString(fmt.Sprintf(`<div class="info-row"><span class="info-k">Modelos</span><span class="info-v">%d</span></div>`, len(s.Program.Models)))
	b.WriteString(fmt.Sprintf(`<div class="info-row"><span class="info-k">Telas</span><span class="info-v">%d</span></div>`, len(s.Program.Screens)))
	b.WriteString(`<div class="info-row"><span class="info-k">Engine</span><span class="info-v">Flang v0.3</span></div>`)
	b.WriteString(`</div></div>`)
	b.WriteString(`</div>`) // dash-grid
	b.WriteString(`</div>`) // dashboard

	// Model sections
	for _, model := range s.Program.Models {
		s.renderSecaoV2(&b, model)
	}

	b.WriteString(`</div></div>`) // content, main

	// Toast
	b.WriteString(`<div id="toast" class="toast"></div>`)

	// JS
	b.WriteString(`<script>`)
	b.WriteString(s.jsV2())
	b.WriteString(`</script></body></html>`)

	return b.String()
}

func (s *Servidor) renderSecaoV2(b *strings.Builder, model *ast.Model) {
	name := lo(model.Name)
	capName := cap(model.Name)

	b.WriteString(fmt.Sprintf(`<div id="secao-%s" class="section" style="display:none">`, name))

	// Header
	b.WriteString(`<div class="sec-head">`)
	b.WriteString(`<div class="sec-left">`)
	b.WriteString(fmt.Sprintf(`<div class="sec-search glass-input"><input type="text" placeholder="Buscar em %s..." oninput="filtrar('%s',this.value)">`, capName, name))
	b.WriteString(svgIcon("search") + `</div></div>`)
	b.WriteString(fmt.Sprintf(`<button class="btn btn-glow" onclick="abrirForm('%s')">`, name))
	b.WriteString(svgIcon("plus") + `<span>Novo ` + capName + `</span></button></div>`)

	// Table card
	b.WriteString(`<div class="card glass-card table-wrap">`)
	b.WriteString(`<table><thead><tr><th class="th-id">#</th>`)
	for _, f := range model.Fields {
		b.WriteString(`<th>` + cap(f.Name) + `</th>`)
	}
	b.WriteString(`<th class="th-act"></th></tr></thead>`)
	b.WriteString(fmt.Sprintf(`<tbody id="tabela-%s"></tbody></table>`, name))
	b.WriteString(fmt.Sprintf(`<div id="vazio-%s" class="empty-state">`, name))
	b.WriteString(svgIcon("inbox") + `<p>Nenhum registro</p></div></div>`)

	// Modal
	b.WriteString(fmt.Sprintf(`<div id="modal-%s" class="modal-wrap" onclick="if(event.target===this)fecharForm('%s')">`, name, name))
	b.WriteString(`<div class="modal glass-modal anim-modal">`)
	b.WriteString(fmt.Sprintf(`<div class="modal-top"><h2 id="titulo-form-%s">Novo %s</h2>`, name, capName))
	b.WriteString(fmt.Sprintf(`<button onclick="fecharForm('%s')" class="modal-x">`, name) + svgIcon("x") + `</button></div>`)
	b.WriteString(fmt.Sprintf(`<form onsubmit="salvar('%s',event)" class="modal-form"><input type="hidden" id="%s-id">`, name, name))

	for _, f := range model.Fields {
		fname := lo(f.Name)
		inputType := tipoInput(f.Type)
		req := ""
		if f.Required {
			req = " required"
		}
		extra := ""
		if f.Type == ast.FieldNumero || f.Type == ast.FieldDinheiro {
			extra = ` step="any"`
		}
		if f.Type == ast.FieldSenha {
			inputType = "password"
		}
		b.WriteString(`<div class="field">`)
		b.WriteString(fmt.Sprintf(`<label for="%s-%s">%s</label>`, name, fname, cap(f.Name)))
		b.WriteString(fmt.Sprintf(`<input type="%s" id="%s-%s" placeholder="%s"%s%s>`,
			inputType, name, fname, cap(f.Name), extra, req))
		b.WriteString(`</div>`)
	}

	b.WriteString(`<div class="modal-foot">`)
	b.WriteString(`<button type="submit" class="btn btn-glow">` + svgIcon("check") + `<span>Salvar</span></button>`)
	b.WriteString(fmt.Sprintf(`<button type="button" class="btn btn-ghost" onclick="fecharForm('%s')">Cancelar</button>`, name))
	b.WriteString(`</div></form></div></div>`)
	b.WriteString(`</div>`)
}

func (s *Servidor) jsV2() string {
	var b strings.Builder

	editIcon := strings.ReplaceAll(strings.ReplaceAll(svgIcon("edit"), `"`, `'`), "\n", "")
	trashIcon := strings.ReplaceAll(strings.ReplaceAll(svgIcon("trash"), `"`, `'`), "\n", "")
	b.WriteString(fmt.Sprintf("var ICO_E=\"%s\",ICO_D=\"%s\";\n", editIcon, trashIcon))

	b.WriteString("var M={\n")
	for _, model := range s.Program.Models {
		name := lo(model.Name)
		b.WriteString(fmt.Sprintf("'%s':[", name))
		for i, f := range model.Fields {
			ft := "t"
			switch f.Type {
			case ast.FieldNumero:
				ft = "n"
			case ast.FieldDinheiro:
				ft = "d"
			case ast.FieldStatus:
				ft = "s"
			case ast.FieldEmail:
				ft = "e"
			}
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(fmt.Sprintf("{n:'%s',t:'%s'}", lo(f.Name), ft))
		}
		b.WriteString("],\n")
	}
	b.WriteString("};\n")

	b.WriteString(`
var ativs=[];
function $(id){return document.getElementById(id);}
function esc(v){if(v==null)return'';var d=document.createElement('div');d.textContent=String(v);return d.innerHTML;}

function irPara(n,el){
  document.querySelectorAll('.section').forEach(function(s){s.style.display='none';});
  var sec=$('secao-'+n);
  if(sec){sec.style.display='block';sec.classList.add('anim-in');}
  document.querySelectorAll('.sb-link').forEach(function(a){a.classList.remove('active');});
  if(el)el.classList.add('active');
  $('page-title').textContent=n==='dashboard'?'Dashboard':n.charAt(0).toUpperCase()+n.slice(1);
  if(innerWidth<768)$('sidebar').classList.remove('open');
}
function irParaNav(n){
  var links=document.querySelectorAll('.sb-link');
  for(var i=0;i<links.length;i++){if(links[i].querySelector('span').textContent.toLowerCase()===n){irPara(n,links[i]);return;}}
  irPara(n,null);
}

function toggleSidebar(){$('sidebar').classList.toggle('open');}
function toggleCollapse(){document.body.classList.toggle('sb-mini');$('sidebar').classList.toggle('mini');}
function toggleDark(){document.body.classList.toggle('dark');}

function toast(msg,t){var e=$('toast');e.textContent=msg;e.className='toast show '+(t||'ok');setTimeout(function(){e.className='toast';},3000);}

function abrirForm(m){$('modal-'+m).classList.add('show');$(m+'-id').value='';$('modal-'+m).querySelector('form').reset();$('titulo-form-'+m).textContent='Novo '+m.charAt(0).toUpperCase()+m.slice(1);}
function fecharForm(m){$('modal-'+m).classList.remove('show');}

function filtrar(m,q){q=q.toLowerCase();document.querySelectorAll('#tabela-'+m+' tr').forEach(function(r){r.style.display=r.textContent.toLowerCase().includes(q)?'':'none';});}

function buscaGlobal(q){
  if(!q){document.querySelectorAll('table tr').forEach(function(r){r.style.display='';});return;}
  q=q.toLowerCase();
  document.querySelectorAll('table tbody tr').forEach(function(r){r.style.display=r.textContent.toLowerCase().includes(q)?'':'none';});
}

function fmtCell(v,t){
  var s=esc(v);
  if(!s||s==='-')return'<span class="muted">—</span>';
  if(t==='s')return'<span class="pill pill-'+pillColor(v)+'">'+s+'</span>';
  if(t==='d'){var n=parseFloat(v);return'<span class="money">R$&nbsp;'+n.toFixed(2)+'</span>';}
  if(t==='e')return'<a class="link" href="mailto:'+s+'">'+s+'</a>';
  return s;
}

function pillColor(v){
  if(!v)return'g';v=v.toLowerCase();
  if('ativo,livre,aberto,ok,sim,disponivel,pronto,entregue,pago,aprovado,online'.indexOf(v)>=0)return'green';
  if('inativo,ocupado,fechado,nao,cancelado,bloqueado,offline,reprovado'.indexOf(v)>=0)return'red';
  if('pendente,aguardando,em andamento,reservado,preparando,analise'.indexOf(v)>=0)return'yellow';
  return'blue';
}

function addAtiv(tipo,mod,nome){
  var labs={c:'Criado',e:'Editado',d:'Excluído'};
  var now=new Date();var h=String(now.getHours()).padStart(2,'0')+':'+String(now.getMinutes()).padStart(2,'0');
  ativs.unshift({t:tipo,m:mod,n:nome,h:h,l:labs[tipo]});
  if(ativs.length>15)ativs.pop();
  renderAtiv();
}
function renderAtiv(){
  var el=$('atividade');
  if(!ativs.length)return;
  var h='';ativs.forEach(function(a){
    h+='<div class="ativ-row"><span class="ativ-tag ativ-'+a.t+'">'+a.l+'</span>';
    h+='<span class="ativ-txt"><b>'+esc(a.m)+'</b>';
    if(a.n)h+=' — '+esc(a.n);
    h+='</span><span class="ativ-time">'+a.h+'</span></div>';
  });
  el.innerHTML=h;
}

function carregar(m){
  fetch('/api/'+m).then(function(r){return r.json();}).then(function(items){
    var tb=$('tabela-'+m),vz=$('vazio-'+m),st=$('stat-'+m);
    tb.innerHTML='';
    if(st)st.textContent=items?items.length:0;
    if(!items||!items.length){vz.style.display='flex';tb.closest('table').style.display='none';return;}
    vz.style.display='none';tb.closest('table').style.display='';
    var cs=M[m];
    items.forEach(function(item){
      var tr=document.createElement('tr');tr.className='row-anim';
      var h='<td class="td-id">'+item.id+'</td>';
      cs.forEach(function(c){h+='<td>'+fmtCell(item[c.n],c.t)+'</td>';});
      h+='<td class="td-act"><button class="act-btn act-edit" onclick="editar(\''+m+'\','+item.id+')">'+ICO_E+'</button>';
      h+='<button class="act-btn act-del" onclick="excluir(\''+m+'\','+item.id+')">'+ICO_D+'</button></td>';
      tr.innerHTML=h;tb.appendChild(tr);
    });
  });
}

function salvar(m,e){
  e.preventDefault();var id=$(m+'-id').value;var d={};
  M[m].forEach(function(c){var v=$(m+'-'+c.n).value;d[c.n]=(c.t==='n'||c.t==='d')?parseFloat(v)||0:v;});
  fetch(id?'/api/'+m+'/'+id:'/api/'+m,{method:id?'PUT':'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(d)})
    .then(function(r){if(!r.ok)return r.json().then(function(e){throw new Error(e.erro);});return r.json();})
    .then(function(){fecharForm(m);carregar(m);addAtiv(id?'e':'c',m,d[M[m][0].n]||'');toast(id?'Atualizado!':'Criado!');})
    .catch(function(err){toast('Erro: '+err.message,'erro');});
}

function editar(m,id){
  fetch('/api/'+m+'/'+id).then(function(r){return r.json();}).then(function(item){
    $(m+'-id').value=item.id;
    M[m].forEach(function(c){$(m+'-'+c.n).value=item[c.n]||'';});
    $('titulo-form-'+m).textContent='Editar';
    $('modal-'+m).classList.add('show');
  });
}

function excluir(m,id){
  if(!confirm('Excluir #'+id+'?'))return;
  var tb=$('tabela-'+m),rows=tb.querySelectorAll('tr'),label='';
  rows.forEach(function(r){if(r.querySelector('.td-id')&&r.querySelector('.td-id').textContent==id){label=r.children[1]?r.children[1].textContent:'';}});
  fetch('/api/'+m+'/'+id,{method:'DELETE'}).then(function(){carregar(m);addAtiv('d',m,label);toast('Excluído!');});
}

// WebSocket - real-time updates
var ws;
function connectWS(){
  var proto=location.protocol==='https:'?'wss:':'ws:';
  ws=new WebSocket(proto+'//'+location.host+'/ws');
  ws.onmessage=function(e){
    try{
      var msg=JSON.parse(e.data);
      if(msg.model) carregar(msg.model);
    }catch(ex){}
  };
  ws.onclose=function(){setTimeout(connectWS,2000);};
  ws.onerror=function(){ws.close();};
}

document.addEventListener('DOMContentLoaded',function(){
  connectWS();
`)
	for _, model := range s.Program.Models {
		b.WriteString(fmt.Sprintf("  carregar('%s');\n", lo(model.Name)))
	}
	b.WriteString("});\n")
	return b.String()
}

// ============================================================
// CSS v2 - 2026 trends: glassmorphism, bento, micro-interactions
// ============================================================

func (s *Servidor) cssV2(theme *ast.Theme) string {
	pri := theme.Primary
	sec := theme.Secondary
	accent := theme.Accent
	side := theme.Sidebar

	return `
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&display=swap');
*{margin:0;padding:0;box-sizing:border-box}

:root{
  --pri:` + pri + `;--sec:` + sec + `;--acc:` + accent + `;--side:` + side + `;
  --bg:#f8fafc;--bg2:#f1f5f9;--card:rgba(255,255,255,.85);--card-solid:#fff;
  --txt:#0f172a;--txt2:#64748b;--txt3:#94a3b8;--brd:rgba(0,0,0,.06);
  --r:16px;--r2:12px;--r3:8px;
  --sh:0 1px 2px rgba(0,0,0,.04),0 2px 8px rgba(0,0,0,.04);
  --sh2:0 4px 24px rgba(0,0,0,.08);
  --sh3:0 8px 40px rgba(0,0,0,.12);
  --ease:cubic-bezier(.4,0,.2,1);--dur:.25s;
  --glass-bg:rgba(255,255,255,.7);--glass-brd:rgba(255,255,255,.4);
  --glass-blur:16px;
}
body.dark{
  --bg:#0c0a1d;--bg2:#12102a;--card:rgba(30,27,75,.6);--card-solid:#1e1b4b;
  --txt:#e2e8f0;--txt2:#94a3b8;--txt3:#64748b;--brd:rgba(255,255,255,.06);
  --sh:0 1px 2px rgba(0,0,0,.2),0 2px 8px rgba(0,0,0,.2);
  --sh2:0 4px 24px rgba(0,0,0,.4);--sh3:0 8px 40px rgba(0,0,0,.5);
  --glass-bg:rgba(15,12,40,.6);--glass-brd:rgba(255,255,255,.08);
}

body{font-family:'Inter',system-ui,-apple-system,sans-serif;background:var(--bg);color:var(--txt);
  display:flex;min-height:100vh;transition:background .4s var(--ease),color .3s var(--ease);
  overflow-x:hidden}

/* ===== Sidebar ===== */
.sidebar{width:260px;background:var(--side);color:#fff;display:flex;flex-direction:column;
  position:fixed;top:0;left:0;bottom:0;z-index:50;transition:width .3s var(--ease),transform .3s var(--ease)}
.sidebar.mini{width:72px}
.sb-top{flex:1;display:flex;flex-direction:column;overflow:hidden}
.sb-brand{padding:20px 16px;display:flex;align-items:center;gap:12px;border-bottom:1px solid rgba(255,255,255,.08)}
.sb-logo{width:36px;height:36px;border-radius:10px;display:flex;align-items:center;justify-content:center;
  background:linear-gradient(135deg,var(--pri),var(--acc));flex-shrink:0}
.sb-logo svg{width:20px;height:20px}
.sb-name{font-weight:700;font-size:1.1rem;white-space:nowrap;overflow:hidden;transition:opacity .2s}
.sb-collapse{margin-left:auto;background:none;border:none;color:rgba(255,255,255,.4);cursor:pointer;padding:4px;
  border-radius:6px;transition:all .2s;flex-shrink:0}
.sb-collapse:hover{color:#fff;background:rgba(255,255,255,.1)}
.sb-collapse svg{width:18px;height:18px;transition:transform .3s}
.sidebar.mini .sb-collapse svg{transform:rotate(180deg)}
.sidebar.mini .sb-name{opacity:0;width:0}
.sidebar.mini .sb-brand{justify-content:center;padding:20px 0}
.sidebar.mini .sb-collapse{display:none}

.sb-nav{padding:12px 8px;display:flex;flex-direction:column;gap:2px;flex:1;overflow-y:auto}
.sb-link{display:flex;align-items:center;gap:12px;padding:10px 12px;border-radius:var(--r3);
  color:rgba(255,255,255,.55);text-decoration:none;font-size:.875rem;font-weight:500;
  transition:all .2s var(--ease);cursor:pointer;white-space:nowrap;position:relative;overflow:hidden}
.sb-link::before{content:'';position:absolute;inset:0;background:rgba(255,255,255,.08);opacity:0;transition:opacity .2s;border-radius:var(--r3)}
.sb-link:hover::before{opacity:1}
.sb-link:hover{color:rgba(255,255,255,.9)}
.sb-icon{width:36px;height:36px;display:flex;align-items:center;justify-content:center;border-radius:var(--r3);
  transition:background .2s;flex-shrink:0}
.sb-icon svg{width:18px;height:18px}
.sb-link.active{color:#fff}
.sb-link.active .sb-icon{background:linear-gradient(135deg,var(--pri),var(--sec));box-shadow:0 2px 12px rgba(99,102,241,.4)}
.sidebar.mini .sb-link span{opacity:0;width:0}
.sidebar.mini .sb-nav{padding:12px 4px}
.sidebar.mini .sb-link{justify-content:center;padding:10px 0}

.sb-foot{padding:12px 16px;border-top:1px solid rgba(255,255,255,.08);display:flex;flex-direction:column;gap:8px}
.sb-theme{display:flex;align-items:center;gap:10px;background:none;border:none;color:rgba(255,255,255,.45);
  cursor:pointer;padding:8px;border-radius:var(--r3);font-size:.85rem;transition:all .2s;width:100%}
.sb-theme:hover{color:#fff;background:rgba(255,255,255,.08)}
.sb-theme svg{width:18px;height:18px;flex-shrink:0}
.sidebar.mini .sb-theme span{display:none}
.sidebar.mini .sb-theme{justify-content:center}
.sidebar.mini .sb-foot{align-items:center}
.sb-powered{font-size:.7rem;color:rgba(255,255,255,.2);text-align:center}
.sidebar.mini .sb-powered{display:none}

/* ===== Main ===== */
.main{margin-left:260px;flex:1;min-height:100vh;transition:margin-left .3s var(--ease)}
body.sb-mini .main{margin-left:72px}

/* ===== Topbar ===== */
.topbar{padding:12px 28px;display:flex;align-items:center;gap:16px;position:sticky;top:0;z-index:30;
  border-bottom:1px solid var(--brd);transition:background .3s}
.topbar h1{font-size:1.1rem;font-weight:700;flex:1;letter-spacing:-.02em}
.tb-menu{display:none;background:none;border:none;color:var(--txt);cursor:pointer;padding:6px;border-radius:var(--r3)}
.tb-menu svg{width:22px;height:22px}
.tb-end{display:flex;align-items:center;gap:12px}
.tb-search{position:relative;display:flex;align-items:center}
.tb-search input{border:none;background:transparent;outline:none;font-size:.875rem;color:var(--txt);width:200px;
  padding:8px 12px 8px 36px;transition:width .3s}
.tb-search input:focus{width:280px}
.tb-search svg{position:absolute;left:10px;width:16px;height:16px;color:var(--txt3);pointer-events:none}

/* ===== Glass ===== */
.glass{background:var(--glass-bg);backdrop-filter:blur(var(--glass-blur));-webkit-backdrop-filter:blur(var(--glass-blur));
  border:1px solid var(--glass-brd)}
.glass-card{background:var(--glass-bg);backdrop-filter:blur(var(--glass-blur));-webkit-backdrop-filter:blur(var(--glass-blur));
  border:1px solid var(--glass-brd);border-radius:var(--r);box-shadow:var(--sh);transition:box-shadow .3s var(--ease),transform .3s var(--ease)}
.glass-card:hover{box-shadow:var(--sh2)}
.glass-input{background:var(--glass-bg);backdrop-filter:blur(8px);border:1px solid var(--glass-brd);border-radius:var(--r2);
  transition:border-color .2s,box-shadow .2s}
.glass-input:focus-within{border-color:var(--pri);box-shadow:0 0 0 3px rgba(99,102,241,.12)}
.glass-modal{background:var(--glass-bg);backdrop-filter:blur(24px);-webkit-backdrop-filter:blur(24px);
  border:1px solid var(--glass-brd);border-radius:var(--r)}

/* ===== Content ===== */
.content{padding:24px 28px}

/* ===== Bento Grid ===== */
.bento{display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:16px;margin-bottom:24px}
.bento-card{position:relative;background:var(--card);border:1px solid var(--brd);border-radius:var(--r);
  padding:24px;cursor:pointer;overflow:hidden;transition:all .3s var(--ease);box-shadow:var(--sh)}
.bento-card:hover{transform:translateY(-4px);box-shadow:var(--sh2);border-color:color-mix(in srgb,var(--accent) 30%,var(--brd))}
.bc-icon{width:48px;height:48px;border-radius:14px;display:flex;align-items:center;justify-content:center;
  background:linear-gradient(135deg,var(--accent),color-mix(in srgb,var(--accent) 70%,#fff));margin-bottom:16px}
.bc-icon svg{width:24px;height:24px;color:#fff}
.bc-num{font-size:clamp(1.75rem,3vw,2.25rem);font-weight:800;letter-spacing:-.03em;line-height:1}
.bc-label{font-size:.85rem;color:var(--txt2);font-weight:500;margin-top:4px}
.bc-glow{position:absolute;top:-40%;right:-20%;width:120px;height:120px;border-radius:50%;
  background:var(--accent);opacity:.06;filter:blur(40px);pointer-events:none;transition:opacity .3s}
.bento-card:hover .bc-glow{opacity:.12}

/* ===== Dashboard Grid ===== */
.dash-grid{display:grid;grid-template-columns:2fr 1fr;gap:16px}
@media(max-width:900px){.dash-grid{grid-template-columns:1fr}}

/* ===== Card ===== */
.card{overflow:hidden}
.card-head{display:flex;align-items:center;gap:10px;padding:18px 20px;border-bottom:1px solid var(--brd)}
.card-head svg{width:18px;height:18px;color:var(--pri)}
.card-head h3{font-size:.95rem;font-weight:600}

/* ===== Activity ===== */
.ativ-list{padding:8px 0;max-height:320px;overflow-y:auto}
.ativ-row{display:flex;align-items:center;gap:10px;padding:10px 20px;font-size:.875rem;transition:background .15s}
.ativ-row:hover{background:rgba(99,102,241,.04)}
.ativ-tag{font-size:.7rem;padding:2px 8px;border-radius:99px;font-weight:700;color:#fff;text-transform:uppercase;letter-spacing:.5px;flex-shrink:0}
.ativ-c{background:#16a34a}.ativ-e{background:var(--pri)}.ativ-d{background:#dc2626}
.ativ-txt{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.ativ-time{color:var(--txt3);font-size:.8rem;font-variant-numeric:tabular-nums;flex-shrink:0}

/* ===== Info ===== */
.info-list{padding:4px 0}
.info-row{display:flex;justify-content:space-between;padding:12px 20px;border-bottom:1px solid var(--brd);font-size:.875rem}
.info-row:last-child{border-bottom:none}
.info-k{color:var(--txt2);font-weight:500}.info-v{font-weight:600}

/* ===== Section ===== */
.section{animation:fadeUp .35s var(--ease)}
@keyframes fadeUp{from{opacity:0;transform:translateY(12px)}to{opacity:1;transform:translateY(0)}}
.anim-in{animation:fadeUp .35s var(--ease)}
.sec-head{display:flex;align-items:center;justify-content:space-between;gap:14px;margin-bottom:20px;flex-wrap:wrap}
.sec-left{flex:1}
.sec-search{display:flex;align-items:center;max-width:380px;padding:0 14px;height:42px}
.sec-search input{flex:1;border:none;background:none;outline:none;font-size:.875rem;color:var(--txt);padding:0 8px}
.sec-search input::placeholder{color:var(--txt3)}
.sec-search svg{width:16px;height:16px;color:var(--txt3);flex-shrink:0}

/* ===== Table ===== */
.table-wrap{overflow-x:auto}
table{width:100%;border-collapse:collapse}
th{text-align:left;padding:12px 16px;font-weight:600;font-size:.75rem;text-transform:uppercase;
  letter-spacing:.6px;color:var(--txt3);background:var(--bg2);border-bottom:1px solid var(--brd)}
td{padding:13px 16px;border-bottom:1px solid var(--brd);font-size:.875rem;transition:background .15s}
tr:hover td{background:rgba(99,102,241,.03)}
.td-id{font-weight:700;color:var(--txt3);font-size:.8rem;width:50px}
.th-id{width:50px}.th-act{width:90px;text-align:right}
.td-act{text-align:right;white-space:nowrap}
.row-anim{animation:fadeUp .25s var(--ease)}

/* Action btns */
.act-btn{width:34px;height:34px;display:inline-flex;align-items:center;justify-content:center;
  border:none;border-radius:var(--r3);cursor:pointer;transition:all .2s var(--ease);background:transparent}
.act-btn svg{width:15px;height:15px}
.act-edit{color:var(--pri)}.act-edit:hover{background:rgba(99,102,241,.1);transform:scale(1.1)}
.act-del{color:#ef4444}.act-del:hover{background:rgba(239,68,68,.1);transform:scale(1.1)}

/* ===== Empty state ===== */
.empty-state{display:flex;flex-direction:column;align-items:center;justify-content:center;padding:48px 20px;color:var(--txt3);gap:8px}
.empty-state svg{width:40px;height:40px;opacity:.4}
.empty-state p{font-size:.9rem}

/* ===== Badges / Pills ===== */
.pill{display:inline-flex;align-items:center;gap:4px;padding:3px 12px;border-radius:99px;font-size:.78rem;font-weight:600;
  text-transform:capitalize;letter-spacing:.2px}
.pill::before{content:'';width:6px;height:6px;border-radius:50%;flex-shrink:0}
.pill-green{background:rgba(22,163,74,.1);color:#16a34a}.pill-green::before{background:#16a34a}
body.dark .pill-green{background:rgba(22,163,74,.15);color:#4ade80}
.pill-red{background:rgba(239,68,68,.1);color:#ef4444}.pill-red::before{background:#ef4444}
body.dark .pill-red{background:rgba(239,68,68,.15);color:#fca5a5}
.pill-yellow{background:rgba(245,158,11,.1);color:#d97706}.pill-yellow::before{background:#f59e0b}
body.dark .pill-yellow{background:rgba(245,158,11,.15);color:#fde047}
.pill-blue{background:rgba(59,130,246,.1);color:#3b82f6}.pill-blue::before{background:#3b82f6}
body.dark .pill-blue{background:rgba(59,130,246,.15);color:#93c5fd}

.money{font-variant-numeric:tabular-nums;font-weight:600;color:var(--pri)}
.muted{color:var(--txt3)}
.link{color:var(--pri);text-decoration:none;font-weight:500}
.link:hover{text-decoration:underline}

/* ===== Buttons ===== */
.btn{display:inline-flex;align-items:center;gap:7px;padding:10px 20px;border:none;border-radius:var(--r2);
  font-size:.875rem;font-weight:600;cursor:pointer;transition:all .25s var(--ease);text-decoration:none;
  position:relative;overflow:hidden}
.btn svg{width:16px;height:16px}
.btn-glow{background:linear-gradient(135deg,var(--pri),var(--sec));color:#fff;
  box-shadow:0 2px 12px rgba(99,102,241,.35)}
.btn-glow:hover{transform:translateY(-2px);box-shadow:0 6px 24px rgba(99,102,241,.45)}
.btn-glow:active{transform:translateY(0)}
.btn-ghost{background:var(--bg2);color:var(--txt2);border:1px solid var(--brd)}
.btn-ghost:hover{background:var(--brd);color:var(--txt)}

/* ===== Modal ===== */
.modal-wrap{display:none;position:fixed;inset:0;background:rgba(0,0,0,.4);backdrop-filter:blur(6px);
  -webkit-backdrop-filter:blur(6px);z-index:100;align-items:center;justify-content:center;padding:20px}
.modal-wrap.show{display:flex}
.modal{width:100%;max-width:460px;max-height:85vh;overflow-y:auto;box-shadow:var(--sh3)}
.anim-modal{animation:modalIn .3s var(--ease)}
@keyframes modalIn{from{opacity:0;transform:scale(.95) translateY(10px)}to{opacity:1;transform:scale(1) translateY(0)}}
.modal-top{display:flex;align-items:center;justify-content:space-between;padding:20px 24px 0}
.modal-top h2{font-size:1.05rem;font-weight:700}
.modal-x{background:none;border:none;color:var(--txt3);cursor:pointer;padding:6px;border-radius:var(--r3);transition:all .2s}
.modal-x:hover{background:var(--bg2);color:var(--txt)}
.modal-x svg{width:18px;height:18px}
.modal-form{padding:16px 24px 24px}

/* ===== Form fields ===== */
.field{margin-bottom:16px}
.field label{display:block;font-weight:600;margin-bottom:6px;font-size:.8rem;color:var(--txt2);text-transform:uppercase;letter-spacing:.5px}
.field input,.field select,.field textarea{width:100%;padding:11px 14px;border:1px solid var(--brd);
  border-radius:var(--r3);font-size:.9rem;background:var(--bg);color:var(--txt);
  transition:all .25s var(--ease);font-family:inherit}
.field input:focus,.field select:focus,.field textarea:focus{outline:none;border-color:var(--pri);
  box-shadow:0 0 0 4px rgba(99,102,241,.1);background:var(--card-solid)}
.field input::placeholder{color:var(--txt3)}
.modal-foot{display:flex;gap:10px;padding-top:16px;border-top:1px solid var(--brd)}

/* ===== Toast ===== */
.toast{position:fixed;bottom:24px;right:24px;padding:14px 28px;border-radius:var(--r2);color:#fff;
  font-weight:600;font-size:.9rem;z-index:200;opacity:0;transform:translateY(16px) scale(.95);
  transition:all .35s var(--ease);pointer-events:none;backdrop-filter:blur(8px)}
.toast.show{opacity:1;transform:translateY(0) scale(1)}
.toast.ok{background:linear-gradient(135deg,#16a34a,#15803d);box-shadow:0 4px 20px rgba(22,163,74,.35)}
.toast.erro{background:linear-gradient(135deg,#ef4444,#dc2626);box-shadow:0 4px 20px rgba(239,68,68,.35)}

/* ===== Scrollbar ===== */
::-webkit-scrollbar{width:6px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--brd);border-radius:3px}
::-webkit-scrollbar-thumb:hover{background:var(--txt3)}

/* ===== Responsive ===== */
@media(max-width:768px){
  .sidebar{transform:translateX(-100%)}
  .sidebar.open{transform:translateX(0)}
  .main{margin-left:0!important}
  .tb-menu{display:block}
  .content{padding:16px}
  .bento{grid-template-columns:1fr 1fr}
  .sec-head{flex-direction:column;align-items:stretch}
  .sec-search{max-width:100%}
  .tb-search input{width:140px}
}
@media(max-width:480px){
  .bento{grid-template-columns:1fr}
}
`
}

// ============================================================
// SVG Icons
// ============================================================

func svgIcon(name string) string {
	icons := map[string]string{
		"zap":      `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>`,
		"layout":   `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>`,
		"menu":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="4" y1="12" x2="20" y2="12"/><line x1="4" y1="6" x2="20" y2="6"/><line x1="4" y1="18" x2="20" y2="18"/></svg>`,
		"search":   `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>`,
		"plus":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>`,
		"edit":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.12 2.12 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>`,
		"trash":    `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>`,
		"x":        `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>`,
		"check":    `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>`,
		"moon":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>`,
		"chevleft": `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>`,
		"activity": `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>`,
		"info":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>`,
		"inbox":    `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 16 12 14 15 10 15 8 12 2 12"/><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/></svg>`,
		"box":      `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>`,
		"users":    `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`,
		"user":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>`,
		"grid":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>`,
		"list":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>`,
		"clip":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1"/></svg>`,
		"dollar":   `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>`,
		"utensils": `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"/><path d="M7 2v20"/><path d="M21 15V2a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3zm0 0v7"/></svg>`,
		"tag":      `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>`,
		"file":     `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>`,
		"settings": `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09a1.65 1.65 0 0 0-1.08-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09a1.65 1.65 0 0 0 1.51-1.08 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9c.26.604.852.997 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>`,
	}
	if svg, ok := icons[name]; ok {
		return svg
	}
	return icons["box"]
}

func modelIcon(name string) string {
	m := map[string]string{
		"produto": "box", "produtos": "box", "prato": "utensils", "pratos": "utensils",
		"cliente": "user", "clientes": "users", "usuario": "user", "usuarios": "users",
		"funcionario": "users", "funcionarios": "users", "equipe": "users",
		"mesa": "grid", "mesas": "grid", "pedido": "clip", "pedidos": "clip",
		"venda": "dollar", "vendas": "dollar", "pagamento": "dollar",
		"categoria": "tag", "categorias": "tag", "item": "list", "itens": "list",
		"configuracao": "settings", "config": "settings", "arquivo": "file",
	}
	if icon, ok := m[name]; ok {
		return icon
	}
	return "box"
}

func tipoInput(ft ast.FieldType) string {
	switch ft {
	case ast.FieldEmail:
		return "email"
	case ast.FieldTelefone:
		return "tel"
	case ast.FieldNumero, ast.FieldDinheiro:
		return "number"
	case ast.FieldData:
		return "date"
	case ast.FieldBooleano:
		return "checkbox"
	case ast.FieldImagem, ast.FieldUpload, ast.FieldArquivo:
		return "file"
	case ast.FieldLink:
		return "url"
	case ast.FieldSenha:
		return "password"
	default:
		return "text"
	}
}

func cap(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
func lo(s string) string { return strings.ToLower(s) }
