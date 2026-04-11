package ide

var ideHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Flang IDE</title>
<script src="https://cdn.tailwindcss.com"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/fabric.js/6.6.1/fabric.min.js"></script>
<script>tailwind.config={darkMode:'class',theme:{extend:{colors:{primary:'#6366f1',accent:'#f59e0b'}}}}</script>
<style>
html,body{margin:0;height:100%;overflow:hidden}
.file-tree{font-size:13px}
.file-item{padding:4px 8px;cursor:pointer;display:flex;align-items:center;gap:6px;border-radius:6px;margin:1px 4px}
.file-item:hover{background:rgba(99,102,241,0.1)}
.file-item.active{background:rgba(99,102,241,0.15);color:#6366f1;font-weight:600}
.file-item.dir{font-weight:500}
.file-children{padding-left:16px}
.tab{padding:6px 16px;font-size:12px;cursor:pointer;border-bottom:2px solid transparent;display:flex;align-items:center;gap:6px;white-space:nowrap}
.tab:hover{background:rgba(255,255,255,0.05)}
.tab.active{border-bottom-color:#6366f1;color:#6366f1;font-weight:600}
.tab .close{opacity:0;font-size:10px;padding:2px 4px;border-radius:4px}
.tab:hover .close{opacity:0.5}
.tab .close:hover{opacity:1;background:rgba(255,255,255,0.1)}
.status-bar{font-size:11px}
.panel-resize{cursor:col-resize;width:4px;background:transparent;transition:background 0.2s}
.panel-resize:hover{background:rgba(99,102,241,0.3)}
#terminal{font-family:'JetBrains Mono','Fira Code',monospace;font-size:12px;line-height:1.6}
#terminal .line{padding:0 12px}
#terminal .error{color:#ef4444}
#terminal .success{color:#22c55e}
#terminal .info{color:#6366f1}
.mode-tab{padding:5px 14px;font-size:12px;border-radius:6px;cursor:pointer;background:none;border:none;color:#94a3b8;transition:all 0.2s}
.mode-tab.active{background:#6366f1;color:white}
.mode-tab:hover:not(.active){background:rgba(99,102,241,0.15);color:#c7d2fe}
.palette-item{display:flex;align-items:center;gap:8px;padding:8px 10px;font-size:12px;border-radius:8px;cursor:grab;background:rgba(99,102,241,0.05);border:1px solid rgba(99,102,241,0.1);color:#c7d2fe;transition:all 0.2s}
.palette-item:hover{background:rgba(99,102,241,0.12);border-color:rgba(99,102,241,0.25)}
.palette-item:active{cursor:grabbing}
.canvas-comp{padding:12px 16px;margin-bottom:8px;border-radius:10px;border:1px solid rgba(255,255,255,0.06);background:rgba(255,255,255,0.03);cursor:pointer;transition:all 0.2s;position:relative}
.canvas-comp:hover{border-color:rgba(99,102,241,0.3);background:rgba(99,102,241,0.05)}
.canvas-comp.selected{border-color:#6366f1;box-shadow:0 0 0 2px rgba(99,102,241,0.2)}
.canvas-comp .comp-delete{position:absolute;top:6px;right:6px;width:20px;height:20px;border-radius:6px;background:rgba(239,68,68,0.1);color:#ef4444;border:none;cursor:pointer;font-size:10px;display:none;align-items:center;justify-content:center}
.canvas-comp:hover .comp-delete{display:flex}
.canvas-comp .comp-label{font-size:10px;text-transform:uppercase;letter-spacing:0.5px;color:#6366f1;font-weight:600;margin-bottom:4px}
.canvas-comp .comp-preview{font-size:13px;color:#94a3b8}
.flow-palette{display:flex;align-items:center;gap:8px;padding:7px 10px;font-size:11px;border-radius:6px;cursor:grab;background:rgba(255,255,255,0.02);color:#c7d2fe;transition:all 0.2s}
.flow-palette:hover{background:rgba(255,255,255,0.05)}
.flow-node{position:absolute;min-width:160px;background:#1e1b4b;border:1px solid rgba(99,102,241,0.3);border-radius:10px;cursor:move;user-select:none;z-index:10}
.flow-node.selected{border-color:#6366f1;box-shadow:0 0 0 2px rgba(99,102,241,0.3)}
.flow-node .node-head{padding:8px 12px;font-size:11px;font-weight:600;border-bottom:1px solid rgba(255,255,255,0.05);display:flex;align-items:center;justify-content:space-between;border-radius:10px 10px 0 0}
.flow-node .node-body{padding:8px 12px;font-size:11px;color:#94a3b8}
.flow-port{width:12px;height:12px;border-radius:50%;border:2px solid #6366f1;background:#0c0a1d;cursor:crosshair;position:absolute;z-index:20}
.flow-port:hover{background:#6366f1}
.flow-port.port-in{top:50%;left:-6px;transform:translateY(-50%)}
.flow-port.port-out{top:50%;right:-6px;transform:translateY(-50%)}
</style>
</head>
<body class="dark bg-gray-950 text-gray-100 flex flex-col h-screen">

<!-- Top bar -->
<div class="flex items-center justify-between px-4 py-2 bg-gray-900 border-b border-gray-800 flex-shrink-0">
  <div class="flex items-center gap-3">
    <div class="w-6 h-6 rounded bg-primary/20 flex items-center justify-center">
      <svg viewBox="0 0 24 24" fill="none" stroke="#6366f1" stroke-width="2" class="w-4 h-4"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>
    </div>
    <span class="font-bold text-sm">Flang IDE</span>
    <span class="text-xs text-gray-500">v0.5.1</span>
    <div class="flex items-center gap-1 bg-gray-800 rounded-lg p-0.5 ml-4">
      <button onclick="switchMode('editor')" id="mode-editor" class="mode-tab active">Codigo</button>
      <button onclick="switchMode('designer')" id="mode-designer" class="mode-tab">Designer</button>
      <button onclick="switchMode('fluxos')" id="mode-fluxos" class="mode-tab">Fluxos</button>
    </div>
  </div>
  <div class="flex items-center gap-2">
    <button onclick="checkProject()" class="px-3 py-1.5 text-xs rounded-lg bg-gray-800 hover:bg-gray-700 text-gray-300 transition-all flex items-center gap-1.5" title="Verificar sintaxe">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-3.5 h-3.5"><polyline points="20 6 9 17 4 12"/></svg>Check
    </button>
    <button onclick="runProject()" class="px-3 py-1.5 text-xs rounded-lg bg-primary hover:bg-primary/80 text-white transition-all flex items-center gap-1.5" title="Executar app">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-3.5 h-3.5"><polygon points="5 3 19 12 5 21 5 3"/></svg>Run
    </button>
    <button onclick="stopProject()" class="px-3 py-1.5 text-xs rounded-lg bg-gray-800 hover:bg-red-500/20 text-gray-400 hover:text-red-400 transition-all flex items-center gap-1.5" title="Parar app">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-3.5 h-3.5"><rect x="6" y="6" width="12" height="12"/></svg>Stop
    </button>
    <a href="http://localhost:8080" target="_blank" class="px-3 py-1.5 text-xs rounded-lg bg-gray-800 hover:bg-gray-700 text-gray-300 transition-all flex items-center gap-1.5" title="Abrir preview">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-3.5 h-3.5"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>Preview
    </a>
  </div>
</div>

<!-- Main layout -->
<div class="flex flex-1 overflow-hidden">

  <!-- File tree sidebar -->
  <div class="w-56 bg-gray-900 border-r border-gray-800 flex flex-col flex-shrink-0 overflow-hidden">
    <div class="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider flex items-center justify-between">
      <span>Arquivos</span>
      <div class="flex gap-1">
        <button onclick="createFile()" class="p-1 rounded hover:bg-gray-800 text-gray-500 hover:text-gray-300" title="Novo arquivo">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-3.5 h-3.5"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        </button>
        <button onclick="loadFileTree()" class="p-1 rounded hover:bg-gray-800 text-gray-500 hover:text-gray-300" title="Atualizar">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-3.5 h-3.5"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
        </button>
      </div>
    </div>
    <div id="file-tree" class="file-tree flex-1 overflow-y-auto px-1"></div>
  </div>

  <!-- Panel: Editor (code) -->
  <div id="panel-editor" class="flex-1 flex flex-col overflow-hidden">
    <!-- Tabs -->
    <div id="tabs" class="flex bg-gray-900 border-b border-gray-800 overflow-x-auto flex-shrink-0"></div>

    <!-- Monaco Editor container -->
    <div id="editor-container" class="flex-1 overflow-hidden"></div>

    <!-- Terminal panel -->
    <div class="border-t border-gray-800 bg-gray-900 flex-shrink-0" style="height:180px">
      <div class="flex items-center justify-between px-3 py-1.5 border-b border-gray-800">
        <span class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Terminal</span>
        <button onclick="clearTerminal()" class="text-xs text-gray-500 hover:text-gray-300">Limpar</button>
      </div>
      <div id="terminal" class="overflow-y-auto p-2" style="height:148px"></div>
    </div>
  </div>

  <!-- Panel: Designer (Fabric.js canvas editor) -->
  <div id="panel-designer" style="display:none" class="flex flex-1 overflow-hidden">

    <!-- Component Palette -->
    <div class="w-52 bg-gray-900 border-r border-gray-800 flex flex-col flex-shrink-0 overflow-hidden">
      <div class="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">Componentes</div>
      <div class="flex-1 overflow-y-auto px-2 space-y-1">
        <div class="palette-item" onmousedown="addToCanvas('titulo')"><span>T</span> Titulo</div>
        <div class="palette-item" onmousedown="addToCanvas('lista')"><span>&#9776;</span> Lista / Tabela</div>
        <div class="palette-item" onmousedown="addToCanvas('botao')"><span>&#9634;</span> Botao</div>
        <div class="palette-item" onmousedown="addToCanvas('busca')"><span>&#128269;</span> Busca</div>
        <div class="palette-item" onmousedown="addToCanvas('grafico')"><span>&#128202;</span> Grafico</div>
        <div class="palette-item" onmousedown="addToCanvas('cards')"><span>&#128203;</span> Cards / Dashboard</div>
        <div class="palette-item" onmousedown="addToCanvas('formulario')"><span>&#128221;</span> Formulario</div>
        <div class="palette-item" onmousedown="addToCanvas('texto')"><span>Aa</span> Texto / Label</div>
        <div class="palette-item" onmousedown="addToCanvas('imagem')"><span>&#128444;</span> Imagem</div>
        <div class="palette-item" onmousedown="addToCanvas('separador')"><span>&#8212;</span> Separador</div>
        <div class="palette-item" onmousedown="addToCanvas('input')"><span>&#9000;</span> Campo de Entrada</div>
        <div class="palette-item" onmousedown="addToCanvas('select')"><span>&#9662;</span> Dropdown</div>
      </div>

      <div class="border-t border-gray-800 p-3">
        <div class="text-xs font-semibold text-gray-500 uppercase mb-2">Tela</div>
        <input type="text" id="d-screen-name" value="principal" placeholder="Nome da tela"
          class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 mb-2 focus:outline-none focus:border-primary"
          oninput="updateDesignerCode()">
        <div class="flex gap-2">
          <button onclick="clearCanvas()" class="flex-1 px-2 py-1.5 text-xs bg-gray-800 hover:bg-gray-700 text-gray-400 rounded-lg transition-all">Limpar</button>
          <button onclick="generateFromDesigner()" class="flex-1 px-2 py-1.5 text-xs bg-primary hover:bg-primary/80 text-white rounded-lg transition-all">Salvar .fg</button>
        </div>
      </div>
    </div>

    <!-- Fabric.js Canvas -->
    <div class="flex-1 relative overflow-hidden bg-gray-950" id="canvas-wrapper">
      <div class="absolute top-3 left-3 z-10 flex items-center gap-2">
        <span id="zoom-level" class="text-xs text-gray-500 bg-gray-900/80 px-2 py-1 rounded">100%</span>
        <button onclick="zoomIn()" class="text-xs text-gray-400 bg-gray-900/80 px-2 py-1 rounded hover:bg-gray-800">+</button>
        <button onclick="zoomOut()" class="text-xs text-gray-400 bg-gray-900/80 px-2 py-1 rounded hover:bg-gray-800">-</button>
        <button onclick="zoomReset()" class="text-xs text-gray-400 bg-gray-900/80 px-2 py-1 rounded hover:bg-gray-800">Reset</button>
      </div>
      <canvas id="fabric-canvas"></canvas>
    </div>

    <!-- Properties + Code -->
    <div class="w-64 bg-gray-900 border-l border-gray-800 flex flex-col flex-shrink-0 overflow-hidden">
      <div class="flex-1 overflow-y-auto p-4">
        <div class="text-xs font-semibold text-gray-500 uppercase mb-3">Propriedades</div>
        <div id="d-props">
          <p class="text-xs text-gray-600">Clique num componente no canvas</p>
        </div>
      </div>
      <div class="border-t border-gray-800 p-4 max-h-[40%] overflow-y-auto">
        <div class="text-xs font-semibold text-gray-500 uppercase mb-2">Codigo .fg Gerado</div>
        <pre id="d-generated" class="text-xs bg-gray-950 text-green-400 p-3 rounded-lg overflow-auto font-mono whitespace-pre-wrap" style="max-height:200px"></pre>
      </div>
    </div>
  </div>

  <!-- Panel: Fluxos (flow logic editor) -->
  <div id="panel-fluxos" style="display:none" class="flex flex-1 overflow-hidden">

    <!-- Node palette -->
    <div class="w-48 bg-gray-900 border-r border-gray-800 p-3 overflow-y-auto flex-shrink-0">
      <div class="text-xs font-semibold text-gray-500 uppercase mb-3">Gatilhos</div>
      <div class="space-y-1">
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'trigger_click')" style="border-left:3px solid #22c55e">Quando Clicar</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'trigger_create')" style="border-left:3px solid #22c55e">Quando Criar</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'trigger_update')" style="border-left:3px solid #22c55e">Quando Atualizar</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'trigger_cron')" style="border-left:3px solid #22c55e">Agendamento</div>
      </div>

      <div class="text-xs font-semibold text-gray-500 uppercase mt-4 mb-3">Acoes</div>
      <div class="space-y-1">
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'action_create')" style="border-left:3px solid #6366f1">Criar Registro</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'action_update')" style="border-left:3px solid #6366f1">Atualizar Registro</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'action_delete')" style="border-left:3px solid #6366f1">Deletar Registro</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'action_message')" style="border-left:3px solid #6366f1">Mostrar Mensagem</div>
      </div>

      <div class="text-xs font-semibold text-gray-500 uppercase mt-4 mb-3">Integracoes</div>
      <div class="space-y-1">
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'integ_whatsapp')" style="border-left:3px solid #22c55e">Enviar WhatsApp</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'integ_email')" style="border-left:3px solid #f59e0b">Enviar Email</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'integ_http')" style="border-left:3px solid #06b6d4">Chamar API</div>
      </div>

      <div class="text-xs font-semibold text-gray-500 uppercase mt-4 mb-3">Logica</div>
      <div class="space-y-1">
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'logic_if')" style="border-left:3px solid #f59e0b">Condicao (Se)</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'logic_loop')" style="border-left:3px solid #f59e0b">Para Cada</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'logic_set')" style="border-left:3px solid #f59e0b">Definir Variavel</div>
        <div class="flow-palette" draggable="true" ondragstart="dragNode(event,'logic_function')" style="border-left:3px solid #f59e0b">Funcao</div>
      </div>

      <div class="mt-6">
        <button onclick="generateFromFlow()" class="w-full px-3 py-2 text-xs bg-primary hover:bg-primary/80 text-white rounded-lg transition-all">Gerar Codigo .fg</button>
        <button onclick="clearFlow()" class="w-full px-3 py-2 text-xs bg-gray-800 text-gray-400 hover:bg-gray-700 rounded-lg transition-all mt-2">Limpar Fluxo</button>
      </div>
    </div>

    <!-- Flow Canvas -->
    <div class="flex-1 relative overflow-hidden bg-gray-950" id="flow-canvas"
         ondragover="event.preventDefault()"
         ondrop="dropNode(event)"
         onmousedown="flowCanvasMouseDown(event)"
         onmousemove="flowCanvasMouseMove(event)"
         onmouseup="flowCanvasMouseUp(event)">
      <svg id="flow-svg" class="absolute inset-0 w-full h-full pointer-events-none" style="z-index:1"></svg>
      <div id="flow-nodes" class="absolute inset-0" style="z-index:2"></div>
      <!-- Grid pattern -->
      <div class="absolute inset-0 opacity-5" style="background-image:radial-gradient(circle,#fff 1px,transparent 1px);background-size:20px 20px"></div>
    </div>

    <!-- Flow Properties -->
    <div class="w-64 bg-gray-900 border-l border-gray-800 p-4 overflow-y-auto flex-shrink-0">
      <div class="text-xs font-semibold text-gray-500 uppercase mb-3">Propriedades do Nodo</div>
      <div id="flow-props">
        <p class="text-xs text-gray-600">Selecione um nodo</p>
      </div>
      <div class="text-xs font-semibold text-gray-500 uppercase mt-6 mb-3">Codigo Gerado</div>
      <pre id="flow-generated-code" class="text-xs bg-gray-950 text-green-400 p-3 rounded-lg overflow-auto max-h-60 font-mono"></pre>
    </div>
  </div>

</div>

<!-- Status bar -->
<div class="status-bar flex items-center justify-between px-4 py-1 bg-primary text-white flex-shrink-0">
  <div class="flex items-center gap-3">
    <span>Flang IDE</span>
    <span id="status-file" class="opacity-70">Nenhum arquivo</span>
  </div>
  <div class="flex items-center gap-3 opacity-70">
    <span id="status-lang">Flang (.fg)</span>
    <span id="status-cursor">Ln 1, Col 1</span>
  </div>
</div>

<!-- Monaco Editor from CDN -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.52.2/min/vs/loader.min.js"></script>
<script>
// State
var editor = null;
var openFiles = {};
var activeFile = null;
var modified = {};

// Mode switching
function switchMode(mode) {
  document.querySelectorAll('.mode-tab').forEach(function(t){t.classList.remove('active');});
  document.getElementById('mode-'+mode).classList.add('active');
  document.getElementById('panel-editor').style.display = mode==='editor' ? 'flex' : 'none';
  document.getElementById('panel-designer').style.display = mode==='designer' ? 'flex' : 'none';
  document.getElementById('panel-fluxos').style.display = mode==='fluxos' ? 'flex' : 'none';
  if (mode === 'designer') {
    setTimeout(function() { initFabricCanvas(); }, 100);
  }
}

// Initialize Monaco
require.config({ paths: { vs: 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.52.2/min/vs' }});
require(['vs/editor/editor.main'], function() {

  // Register Flang language
  monaco.languages.register({ id: 'flang', extensions: ['.fg'] });

  // Syntax highlighting
  monaco.languages.setMonarchTokensProvider('flang', {
    keywords: ['sistema','dados','telas','tela','eventos','tema','logica','banco','autenticacao','integracoes',
      'importar','de','system','models','screens','screen','events','theme','logic','database','auth','import','from',
      'rotas','rota','paginas','pagina','sidebar','item'],
    typeKeywords: ['texto','texto_longo','numero','dinheiro','email','telefone','data','booleano','imagem','arquivo',
      'upload','link','status','senha','enum','text','number','money','phone','image','file','password','boolean','date'],
    controlKeywords: ['se','senao','enquanto','repetir','para','para_cada','funcao','retornar','definir','mostrar',
      'quando','clicar','criar','atualizar','deletar','enviar','validar','tentar','erro','parar','continuar',
      'if','else','while','repeat','for','function','return','set','print','when','click','create','update','delete','try','error','break','continue'],
    modifiers: ['obrigatorio','unico','pertence_a','tem_muitos','muitos_para_muitos','soft_delete','indice','padrao',
      'required','unique','belongs_to','has_many','many_to_many','index','default'],
    screenKw: ['titulo','lista','mostrar','botao','busca','grafico','dashboard','formulario','tabela','requer','publico',
      'title','list','show','button','search','chart','form','table','requires','public'],
    themeKw: ['cor','primaria','secundaria','destaque','escuro','claro','fonte','borda','fundo','estilo','icone',
      'moderno','simples','elegante','corporativo','glassmorphism','flat','neumorphism','minimal'],
    colors: ['azul','verde','vermelho','roxo','laranja','rosa','amarelo','ciano','indigo','branco','preto'],
    operators: ['==','!=','>=','<=','>','<','+','-','*','/','='],
    tokenizer: {
      root: [
        [/#.*$/, 'comment'],
        [/\/\/.*$/, 'comment'],
        [/"[^"]*"/, 'string'],
        [/\b\d+(\.\d+)?\b/, 'number'],
        [/\b(verdadeiro|falso|nulo|nada|true|false|null)\b/, 'constant'],
        [/[a-zA-Z_\u00C0-\u024F\u0400-\u04FF\u4E00-\u9FFF\u3040-\u309F\u30A0-\u30FF\uAC00-\uD7AF\u0600-\u06FF\u0900-\u097F\u0980-\u09FF\u0E00-\u0E7F][a-zA-Z0-9_\u00C0-\u024F\u0400-\u04FF\u4E00-\u9FFF\u3040-\u309F\u30A0-\u30FF\uAC00-\uD7AF\u0600-\u06FF\u0900-\u097F\u0980-\u09FF\u0E00-\u0E7F]*/, {
          cases: {
            '@keywords': 'keyword',
            '@typeKeywords': 'type',
            '@controlKeywords': 'keyword.control',
            '@modifiers': 'keyword.modifier',
            '@screenKw': 'keyword.screen',
            '@themeKw': 'keyword.theme',
            '@colors': 'constant.color',
            '@default': 'identifier'
          }
        }],
        [/[{}()\[\]]/, 'delimiter'],
        [/[;,.]/, 'delimiter'],
        [/:/, 'delimiter.colon'],
      ]
    }
  });

  // Custom theme
  monaco.editor.defineTheme('flang-dark', {
    base: 'vs-dark',
    inherit: true,
    rules: [
      { token: 'keyword', foreground: '6366f1', fontStyle: 'bold' },
      { token: 'keyword.control', foreground: 'c084fc' },
      { token: 'keyword.modifier', foreground: 'f59e0b' },
      { token: 'keyword.screen', foreground: '22c55e' },
      { token: 'keyword.theme', foreground: 'ec4899' },
      { token: 'type', foreground: '06b6d4', fontStyle: 'italic' },
      { token: 'string', foreground: 'a5f3fc' },
      { token: 'number', foreground: 'fbbf24' },
      { token: 'comment', foreground: '475569', fontStyle: 'italic' },
      { token: 'constant', foreground: 'f97316' },
      { token: 'constant.color', foreground: '34d399' },
      { token: 'delimiter.colon', foreground: '94a3b8' },
    ],
    colors: {
      'editor.background': '#0c0a1d',
      'editor.foreground': '#e2e8f0',
      'editor.lineHighlightBackground': '#1e1b4b30',
      'editor.selectionBackground': '#6366f140',
      'editorCursor.foreground': '#6366f1',
      'editorLineNumber.foreground': '#334155',
      'editorLineNumber.activeForeground': '#6366f1',
    }
  });

  // Create editor
  editor = monaco.editor.create(document.getElementById('editor-container'), {
    value: '# Bem-vindo ao Flang IDE!\n# Selecione um arquivo na arvore a esquerda.\n# Ou clique + para criar um novo arquivo .fg\n',
    language: 'flang',
    theme: 'flang-dark',
    fontSize: 14,
    fontFamily: "'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace",
    minimap: { enabled: true, scale: 1 },
    smoothScrolling: true,
    cursorBlinking: 'smooth',
    cursorSmoothCaretAnimation: 'on',
    padding: { top: 12, bottom: 12 },
    renderLineHighlight: 'all',
    bracketPairColorization: { enabled: true },
    automaticLayout: true,
    wordWrap: 'on',
    tabSize: 2,
    scrollBeyondLastLine: false,
  });

  // Track cursor position
  editor.onDidChangeCursorPosition(function(e) {
    document.getElementById('status-cursor').textContent = 'Ln ' + e.position.lineNumber + ', Col ' + e.position.column;
  });

  // Track modifications
  editor.onDidChangeModelContent(function() {
    if (activeFile) {
      modified[activeFile] = true;
      updateTabModified(activeFile, true);
    }
  });

  // Keyboard shortcuts
  editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, function() {
    saveCurrentFile();
  });

  // Load file tree
  loadFileTree();

  termLog('info', 'Flang IDE iniciado. Pronto para editar!');
});

// File tree
function loadFileTree() {
  fetch('/api/files').then(function(r){return r.json();}).then(function(files) {
    document.getElementById('file-tree').innerHTML = renderTree(files);
  });
}

function renderTree(files) {
  if (!files || !files.length) return '<div class="px-3 py-4 text-xs text-gray-600">Nenhum arquivo</div>';
  var html = '';
  files.forEach(function(f) {
    if (f.isDir) {
      html += '<div class="file-item dir" onclick="this.nextElementSibling.classList.toggle(\'hidden\')">'+
        '<svg viewBox="0 0 24 24" fill="none" stroke="#f59e0b" stroke-width="2" class="w-4 h-4 flex-shrink-0"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>'+
        '<span>'+f.name+'</span></div>';
      html += '<div class="file-children">' + renderTree(f.children) + '</div>';
    } else {
      var icon = f.name.endsWith('.fg') ?
        '<svg viewBox="0 0 24 24" fill="none" stroke="#6366f1" stroke-width="2" class="w-4 h-4 flex-shrink-0"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>' :
        '<svg viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" class="w-4 h-4 flex-shrink-0"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>';
      html += '<div class="file-item" onclick="openFile(\''+f.path+'\')" data-path="'+f.path+'">'+icon+'<span>'+f.name+'</span></div>';
    }
  });
  return html;
}

// Open file
function openFile(path) {
  // Highlight in tree
  document.querySelectorAll('.file-item').forEach(function(el){el.classList.remove('active');});
  var el = document.querySelector('[data-path="'+path+'"]');
  if(el) el.classList.add('active');

  if (openFiles[path]) {
    switchToFile(path);
    return;
  }

  fetch('/api/file?path='+encodeURIComponent(path)).then(function(r){return r.text();}).then(function(content) {
    openFiles[path] = content;
    addTab(path);
    switchToFile(path);
    document.getElementById('status-file').textContent = path;
  });
}

function switchToFile(path) {
  activeFile = path;
  var lang = path.endsWith('.fg') ? 'flang' : (path.endsWith('.json') ? 'json' : (path.endsWith('.go') ? 'go' : (path.endsWith('.js') ? 'javascript' : 'plaintext')));
  editor.setValue(openFiles[path] || '');
  monaco.editor.setModelLanguage(editor.getModel(), lang);

  // Update tabs
  document.querySelectorAll('.tab').forEach(function(t){t.classList.remove('active');});
  var tab = document.querySelector('.tab[data-path="'+path+'"]');
  if(tab) tab.classList.add('active');

  // Update tree highlight
  document.querySelectorAll('.file-item').forEach(function(el){el.classList.remove('active');});
  var el = document.querySelector('[data-path="'+path+'"]');
  if(el) el.classList.add('active');

  document.getElementById('status-file').textContent = path;
  document.getElementById('status-lang').textContent = lang === 'flang' ? 'Flang (.fg)' : lang;
}

// Tabs
function addTab(path) {
  var name = path.split('/').pop();
  var tabs = document.getElementById('tabs');
  if (document.querySelector('.tab[data-path="'+path+'"]')) return;

  var tab = document.createElement('div');
  tab.className = 'tab active';
  tab.setAttribute('data-path', path);
  tab.innerHTML = '<span onclick="switchToFile(\''+path+'\')">'+name+'</span><span class="close" onclick="event.stopPropagation();closeTab(\''+path+'\')">&times;</span>';
  tab.onclick = function(){switchToFile(path);};
  tabs.appendChild(tab);

  document.querySelectorAll('.tab').forEach(function(t){t.classList.remove('active');});
  tab.classList.add('active');
}

function closeTab(path) {
  var tab = document.querySelector('.tab[data-path="'+path+'"]');
  if(tab) tab.remove();
  delete openFiles[path];
  delete modified[path];

  var remaining = document.querySelectorAll('.tab');
  if (remaining.length > 0) {
    var last = remaining[remaining.length-1];
    switchToFile(last.getAttribute('data-path'));
  } else {
    activeFile = null;
    editor.setValue('# Selecione um arquivo');
    document.getElementById('status-file').textContent = 'Nenhum arquivo';
  }
}

function updateTabModified(path, isModified) {
  var tab = document.querySelector('.tab[data-path="'+path+'"] span:first-child');
  if (!tab) return;
  var name = path.split('/').pop();
  tab.textContent = isModified ? name + ' \u25cf' : name;
}

// Save
function saveCurrentFile() {
  if (!activeFile) return;
  var content = editor.getValue();
  openFiles[activeFile] = content;

  fetch('/api/file/save', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({path: activeFile, content: content})
  }).then(function(r) {
    if (r.ok) {
      modified[activeFile] = false;
      updateTabModified(activeFile, false);
      termLog('success', 'Salvo: ' + activeFile);
    }
  });
}

// Create file
function createFile() {
  var name = prompt('Nome do arquivo (ex: dados/produto.fg):');
  if (!name) return;
  fetch('/api/file/create', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({path: name, isDir: false})
  }).then(function() {
    loadFileTree();
    setTimeout(function(){openFile(name);}, 500);
  });
}

// Run/Check/Stop
function runProject() {
  var file = findMainFile();
  termLog('info', 'Executando ' + file + '...');
  fetch('/api/run', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({file: file})
  }).then(function(r){return r.json();}).then(function(d) {
    if (d.status === 'running') {
      termLog('success', 'App rodando em ' + d.url);
    } else {
      termLog('error', 'Erro: ' + (d.message||'desconhecido'));
    }
  });
}

function stopProject() {
  fetch('/api/stop').then(function(r){return r.json();}).then(function() {
    termLog('info', 'App parado.');
  });
}

function checkProject() {
  var file = findMainFile();
  termLog('info', 'Verificando ' + file + '...');
  fetch('/api/check', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({file: file})
  }).then(function(r){return r.json();}).then(function(d) {
    var lines = d.output.split('\n');
    lines.forEach(function(line) {
      if (!line.trim()) return;
      if (line.indexOf('ERRO') >= 0) termLog('error', line);
      else if (line.indexOf('valido') >= 0 || line.indexOf('OK') >= 0) termLog('success', line);
      else termLog('info', line);
    });
  });
}

function findMainFile() {
  if (activeFile && activeFile.endsWith('.fg')) return activeFile;
  return 'inicio.fg';
}

// Terminal
function termLog(type, msg) {
  var term = document.getElementById('terminal');
  var time = new Date().toLocaleTimeString();
  term.innerHTML += '<div class="line ' + type + '"><span style="opacity:0.4">['+time+']</span> '+msg+'</div>';
  term.scrollTop = term.scrollHeight;
}

function clearTerminal() {
  document.getElementById('terminal').innerHTML = '';
}

// ============================================================
// DESIGNER - Fabric.js Canvas Editor
// ============================================================

var fabricCanvas = null;
var designerReady = false;
var compCounter = 0;

function initFabricCanvas() {
  if (designerReady) return;

  var wrapper = document.getElementById('canvas-wrapper');
  var cEl = document.getElementById('fabric-canvas');
  cEl.width = wrapper.clientWidth;
  cEl.height = wrapper.clientHeight;

  fabricCanvas = new fabric.Canvas('fabric-canvas', {
    backgroundColor: '#0a0a1a',
    selection: true,
    preserveObjectStacking: true,
  });

  // Grid background via CSS
  wrapper.style.backgroundImage = 'radial-gradient(circle, rgba(99,102,241,0.08) 1px, transparent 1px)';
  wrapper.style.backgroundSize = '20px 20px';

  // Zoom with mouse wheel
  fabricCanvas.on('mouse:wheel', function(opt) {
    var delta = opt.e.deltaY;
    var zoom = fabricCanvas.getZoom();
    zoom *= Math.pow(0.999, delta);
    if (zoom > 5) zoom = 5;
    if (zoom < 0.2) zoom = 0.2;
    fabricCanvas.zoomToPoint({ x: opt.e.offsetX, y: opt.e.offsetY }, zoom);
    document.getElementById('zoom-level').textContent = Math.round(zoom * 100) + '%';
    opt.e.preventDefault();
    opt.e.stopPropagation();
  });

  // Pan with Alt+drag or middle mouse
  var panning = false;
  fabricCanvas.on('mouse:down', function(opt) {
    if (opt.e.altKey || opt.e.button === 1) {
      panning = true;
      fabricCanvas.selection = false;
      fabricCanvas.setCursor('grab');
    }
  });
  fabricCanvas.on('mouse:move', function(opt) {
    if (panning) {
      var e = opt.e;
      var vpt = fabricCanvas.viewportTransform;
      vpt[4] += e.movementX;
      vpt[5] += e.movementY;
      fabricCanvas.requestRenderAll();
    }
  });
  fabricCanvas.on('mouse:up', function() {
    panning = false;
    fabricCanvas.selection = true;
    fabricCanvas.setCursor('default');
  });

  // Select/deselect
  fabricCanvas.on('selection:created', function(e) { showDesignerProps(e.selected[0]); });
  fabricCanvas.on('selection:updated', function(e) { showDesignerProps(e.selected[0]); });
  fabricCanvas.on('selection:cleared', function() { clearDesignerProps(); });

  // Update code on move/resize
  fabricCanvas.on('object:modified', function() { updateDesignerCode(); });
  fabricCanvas.on('object:added', function() { updateDesignerCode(); });
  fabricCanvas.on('object:removed', function() { updateDesignerCode(); });

  // Delete key
  document.addEventListener('keydown', function(e) {
    if (e.key === 'Delete' && fabricCanvas && fabricCanvas.getActiveObject()) {
      var obj = fabricCanvas.getActiveObject();
      if (obj && obj.compType) {
        fabricCanvas.remove(obj);
        fabricCanvas.discardActiveObject();
        clearDesignerProps();
        updateDesignerCode();
      }
    }
  });

  // Resize observer
  new ResizeObserver(function() {
    if (!fabricCanvas) return;
    fabricCanvas.setDimensions({
      width: wrapper.clientWidth,
      height: wrapper.clientHeight
    });
    fabricCanvas.renderAll();
  }).observe(wrapper);

  designerReady = true;
}

function addToCanvas(type) {
  if (!fabricCanvas) initFabricCanvas();

  compCounter++;
  var id = 'comp-' + compCounter;
  var group;

  var cx = fabricCanvas.width / 2 / fabricCanvas.getZoom();
  var cy = fabricCanvas.height / 2 / fabricCanvas.getZoom();
  cx += (Math.random() - 0.5) * 100;
  cy += (Math.random() - 0.5) * 100;

  switch(type) {
    case 'titulo':
      group = createCompGroup(id, type, cx - 150, cy - 20, 300, 50, '#6366f1', 'Titulo', {texto: 'Minha Tela'});
      break;
    case 'lista':
      group = createCompGroup(id, type, cx - 200, cy - 80, 400, 180, '#3b82f6', 'Lista / Tabela', {modelo: 'produto', campos: 'nome, preco, status'});
      break;
    case 'botao':
      group = createCompGroup(id, type, cx - 60, cy - 18, 140, 40, '#22c55e', 'Botao', {texto: 'Novo', cor: 'azul', acao: 'criar produto'});
      break;
    case 'busca':
      group = createCompGroup(id, type, cx - 150, cy - 18, 300, 40, '#06b6d4', 'Busca', {modelo: 'produto'});
      break;
    case 'grafico':
      group = createCompGroup(id, type, cx - 180, cy - 80, 360, 180, '#f59e0b', 'Grafico', {modelo: 'produto', tipo: 'barra'});
      break;
    case 'cards':
      group = createCompGroup(id, type, cx - 200, cy - 50, 400, 100, '#8b5cf6', 'Dashboard Cards', {});
      break;
    case 'formulario':
      group = createCompGroup(id, type, cx - 160, cy - 100, 320, 220, '#ec4899', 'Formulario', {modelo: 'produto'});
      break;
    case 'texto':
      group = createCompGroup(id, type, cx - 100, cy - 15, 200, 35, '#64748b', 'Texto', {conteudo: 'Texto aqui...'});
      break;
    case 'imagem':
      group = createCompGroup(id, type, cx - 80, cy - 80, 160, 160, '#14b8a6', 'Imagem', {url: ''});
      break;
    case 'separador':
      group = createCompGroup(id, type, cx - 150, cy - 2, 300, 4, '#475569', 'Separador', {});
      break;
    case 'input':
      group = createCompGroup(id, type, cx - 120, cy - 18, 240, 40, '#a855f7', 'Campo', {nome: 'campo', tipo: 'texto'});
      break;
    case 'select':
      group = createCompGroup(id, type, cx - 120, cy - 18, 240, 40, '#f97316', 'Dropdown', {nome: 'campo', opcoes: 'opcao1, opcao2, opcao3'});
      break;
  }

  if (group) {
    fabricCanvas.add(group);
    fabricCanvas.setActiveObject(group);
    fabricCanvas.renderAll();
    showDesignerProps(group);
  }
}

function createCompGroup(id, type, x, y, w, h, color, label, props) {
  var bg = new fabric.Rect({
    width: w, height: h,
    fill: color + '10',
    stroke: color + '40',
    strokeWidth: 1.5,
    rx: 8, ry: 8,
    originX: 'left', originY: 'top',
  });

  var labelText = new fabric.Text(label, {
    fontSize: 10,
    fill: color,
    fontFamily: 'Inter, system-ui, sans-serif',
    fontWeight: '700',
    left: 8,
    top: 6,
    originX: 'left', originY: 'top',
  });

  var previewText = '';
  switch(type) {
    case 'titulo': previewText = props.texto || 'Titulo'; break;
    case 'lista': previewText = 'ID | ' + (props.campos||'nome, preco').split(',').join(' | '); break;
    case 'botao': previewText = '[ ' + (props.texto||'Botao') + ' ]'; break;
    case 'busca': previewText = 'Buscar em ' + (props.modelo||'...'); break;
    case 'grafico': previewText = (props.tipo||'barra') + ' - ' + (props.modelo||'dados'); break;
    case 'cards': previewText = '[ Card 1 ] [ Card 2 ] [ Card 3 ]'; break;
    case 'formulario': previewText = 'Nome: [____]  Email: [____]  [ Salvar ]'; break;
    case 'texto': previewText = props.conteudo || 'Texto'; break;
    case 'imagem': previewText = '[ Imagem ]'; break;
    case 'separador': previewText = ''; break;
    case 'input': previewText = (props.nome||'campo') + ': [____________]'; break;
    case 'select': previewText = (props.nome||'campo') + ': [v selecione ]'; break;
  }

  var content = new fabric.Text(previewText, {
    fontSize: 12,
    fill: '#94a3b8',
    fontFamily: "'JetBrains Mono', monospace",
    left: 8,
    top: 22,
    originX: 'left', originY: 'top',
  });

  var group = new fabric.Group([bg, labelText, content], {
    left: x, top: y,
    originX: 'left', originY: 'top',
    cornerStyle: 'circle',
    cornerColor: color,
    cornerStrokeColor: '#fff',
    cornerSize: 8,
    transparentCorners: false,
    borderColor: color,
    borderScaleFactor: 2,
    padding: 4,
    snapAngle: 45,
  });

  group.compId = id;
  group.compType = type;
  group.compProps = props;
  group.compColor = color;
  group.compLabel = label;

  return group;
}

function showDesignerProps(obj) {
  if (!obj || !obj.compType) { clearDesignerProps(); return; }
  var panel = document.getElementById('d-props');
  var type = obj.compType;
  var props = obj.compProps || {};
  var color = obj.compColor || '#6366f1';

  var fields = {
    titulo: [{key:'texto',label:'Texto',type:'text'}],
    lista: [{key:'modelo',label:'Modelo',type:'text'},{key:'campos',label:'Campos (virgula)',type:'text'}],
    botao: [{key:'texto',label:'Texto',type:'text'},{key:'cor',label:'Cor',type:'select',options:['azul','verde','vermelho','amarelo','roxo']},{key:'acao',label:'Acao',type:'text'}],
    busca: [{key:'modelo',label:'Modelo',type:'text'}],
    grafico: [{key:'modelo',label:'Modelo',type:'text'},{key:'tipo',label:'Tipo',type:'select',options:['barra','pizza','linha','doughnut']}],
    cards: [],
    formulario: [{key:'modelo',label:'Modelo',type:'text'}],
    texto: [{key:'conteudo',label:'Conteudo',type:'textarea'}],
    imagem: [{key:'url',label:'URL da imagem',type:'text'}],
    separador: [],
    input: [{key:'nome',label:'Nome do campo',type:'text'},{key:'tipo',label:'Tipo',type:'select',options:['texto','numero','email','telefone','data','senha','dinheiro']}],
    select: [{key:'nome',label:'Nome do campo',type:'text'},{key:'opcoes',label:'Opcoes (virgula)',type:'text'}]
  };

  var html = '<div class="space-y-3">';
  html += '<div class="flex items-center gap-2 mb-2"><div class="w-3 h-3 rounded" style="background:'+color+'"></div><span class="text-xs font-bold" style="color:'+color+'">'+obj.compLabel+'</span></div>';

  (fields[type]||[]).forEach(function(f) {
    html += '<div><label class="block text-xs text-gray-400 mb-1">'+f.label+'</label>';
    if (f.type === 'select') {
      html += '<select onchange="updateCanvasProp(\''+f.key+'\',this.value)" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 focus:outline-none focus:border-primary">';
      f.options.forEach(function(o) {
        html += '<option value="'+o+'"'+(props[f.key]===o?' selected':'')+'>'+o+'</option>';
      });
      html += '</select>';
    } else if (f.type === 'textarea') {
      html += '<textarea onchange="updateCanvasProp(\''+f.key+'\',this.value)" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 focus:outline-none focus:border-primary resize-none" rows="3">'+(props[f.key]||'')+'</textarea>';
    } else {
      html += '<input type="text" value="'+(props[f.key]||'')+'" onchange="updateCanvasProp(\''+f.key+'\',this.value)" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 focus:outline-none focus:border-primary">';
    }
    html += '</div>';
  });

  html += '<div class="pt-3 border-t border-gray-800 mt-3">';
  html += '<div class="grid grid-cols-2 gap-2 text-xs text-gray-500">';
  html += '<span>X: '+Math.round(obj.left)+'</span><span>Y: '+Math.round(obj.top)+'</span>';
  html += '<span>W: '+Math.round(obj.width * obj.scaleX)+'</span><span>H: '+Math.round(obj.height * obj.scaleY)+'</span>';
  html += '</div></div>';

  html += '<button onclick="deleteSelected()" class="w-full mt-3 px-3 py-1.5 text-xs bg-red-500/10 text-red-400 hover:bg-red-500/20 rounded-lg transition-all">Remover componente</button>';
  html += '</div>';
  panel.innerHTML = html;
}

function clearDesignerProps() {
  document.getElementById('d-props').innerHTML = '<p class="text-xs text-gray-600">Clique num componente no canvas</p>';
}

function updateCanvasProp(key, val) {
  var obj = fabricCanvas.getActiveObject();
  if (!obj || !obj.compProps) return;
  obj.compProps[key] = val;

  var items = obj.getObjects();
  if (items.length >= 3) {
    var previewText = '';
    switch(obj.compType) {
      case 'titulo': previewText = obj.compProps.texto || 'Titulo'; break;
      case 'lista': previewText = 'ID | ' + (obj.compProps.campos||'').split(',').join(' | '); break;
      case 'botao': previewText = '[ ' + (obj.compProps.texto||'Botao') + ' ]'; break;
      case 'busca': previewText = 'Buscar em ' + (obj.compProps.modelo||'...'); break;
      case 'grafico': previewText = (obj.compProps.tipo||'barra') + ' - ' + (obj.compProps.modelo||'dados'); break;
      case 'formulario': previewText = 'Formulario: ' + (obj.compProps.modelo||'item'); break;
      case 'texto': previewText = obj.compProps.conteudo || 'Texto'; break;
      case 'input': previewText = (obj.compProps.nome||'campo') + ': [____________]'; break;
      case 'select': previewText = (obj.compProps.nome||'campo') + ': [v selecione ]'; break;
      default: previewText = obj.compLabel;
    }
    items[2].set('text', previewText);
  }
  fabricCanvas.renderAll();
  updateDesignerCode();
}

function deleteSelected() {
  var obj = fabricCanvas.getActiveObject();
  if (obj) {
    fabricCanvas.remove(obj);
    fabricCanvas.discardActiveObject();
    clearDesignerProps();
    updateDesignerCode();
  }
}

function clearCanvas() {
  if (!fabricCanvas) return;
  var objects = fabricCanvas.getObjects().filter(function(o){ return o.compType; });
  objects.forEach(function(o){ fabricCanvas.remove(o); });
  fabricCanvas.discardActiveObject();
  clearDesignerProps();
  updateDesignerCode();
}

function zoomIn() {
  if (!fabricCanvas) return;
  var z = fabricCanvas.getZoom() * 1.2;
  if (z > 5) z = 5;
  fabricCanvas.setZoom(z);
  document.getElementById('zoom-level').textContent = Math.round(z*100)+'%';
}
function zoomOut() {
  if (!fabricCanvas) return;
  var z = fabricCanvas.getZoom() * 0.8;
  if (z < 0.2) z = 0.2;
  fabricCanvas.setZoom(z);
  document.getElementById('zoom-level').textContent = Math.round(z*100)+'%';
}
function zoomReset() {
  if (!fabricCanvas) return;
  fabricCanvas.setZoom(1);
  fabricCanvas.viewportTransform = [1,0,0,1,0,0];
  fabricCanvas.renderAll();
  document.getElementById('zoom-level').textContent = '100%';
}

function updateDesignerCode() {
  if (!fabricCanvas) return;
  var objects = fabricCanvas.getObjects().filter(function(o){ return o.compType; });
  objects.sort(function(a,b){ return a.top - b.top; });

  var screenName = document.getElementById('d-screen-name').value || 'principal';
  var code = 'telas\n\n';
  code += '  tela ' + screenName + '\n';

  objects.forEach(function(obj) {
    var p = obj.compProps || {};
    switch(obj.compType) {
      case 'titulo':
        code += '    titulo "' + (p.texto||'Titulo') + '"\n';
        break;
      case 'lista':
        code += '    lista ' + (p.modelo||'item') + '\n';
        (p.campos||'').split(',').forEach(function(f) {
          f = f.trim();
          if (f) code += '      mostrar ' + f + '\n';
        });
        break;
      case 'botao':
        code += '    botao ' + (p.cor||'azul') + '\n';
        code += '      texto "' + (p.texto||'Novo') + '"\n';
        break;
      case 'busca':
        code += '    busca ' + (p.modelo||'item') + '\n';
        break;
      case 'grafico':
        code += '    grafico ' + (p.modelo||'item') + '\n';
        code += '      tipo ' + (p.tipo||'barra') + '\n';
        break;
      case 'cards':
        code += '    dashboard\n';
        break;
      case 'formulario':
        code += '    formulario ' + (p.modelo||'item') + '\n';
        break;
      case 'texto':
        code += '    # ' + (p.conteudo||'Texto') + '\n';
        break;
      case 'imagem':
        if (p.url) code += '    # imagem: ' + p.url + '\n';
        break;
      case 'separador':
        code += '    # ---\n';
        break;
      case 'input':
        code += '    # campo ' + (p.nome||'campo') + ': ' + (p.tipo||'texto') + '\n';
        break;
      case 'select':
        code += '    # select ' + (p.nome||'campo') + ': ' + (p.opcoes||'') + '\n';
        break;
    }
  });

  code += '\n';
  document.getElementById('d-generated').textContent = code;
}

function generateFromDesigner() {
  updateDesignerCode();
  var code = document.getElementById('d-generated').textContent;
  if (!code || code.trim() === 'telas\n\n  tela principal') {
    termLog('error', 'Canvas vazio - arraste componentes primeiro');
    return;
  }
  var filename = 'telas_visual.fg';
  fetch('/api/file/save', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({path: filename, content: code})
  }).then(function() {
    termLog('success', 'Telas salvas em ' + filename);
    loadFileTree();
    openFile(filename);
    switchMode('editor');
  });
}

// ============================================================
// FLUXOS - Flow Logic Editor
// ============================================================

var flowNodes = [];
var flowConnections = [];
var flowNodeId = 0;
var selectedNode = null;
var draggingNode = null;
var dragOffset = {x:0, y:0};
var connecting = false;
var connectFrom = null;

var nodeTypes = {
  trigger_click: {label:'Quando Clicar', color:'#22c55e', category:'trigger', fields:[{key:'botao',label:'Botao',type:'text'}]},
  trigger_create: {label:'Quando Criar', color:'#22c55e', category:'trigger', fields:[{key:'modelo',label:'Modelo',type:'text'}]},
  trigger_update: {label:'Quando Atualizar', color:'#22c55e', category:'trigger', fields:[{key:'modelo',label:'Modelo',type:'text'}]},
  trigger_cron: {label:'Agendamento', color:'#22c55e', category:'trigger', fields:[{key:'intervalo',label:'Intervalo',type:'text'},{key:'unidade',label:'Unidade',type:'select',options:['minutos','horas']}]},
  action_create: {label:'Criar Registro', color:'#6366f1', category:'action', fields:[{key:'modelo',label:'Modelo',type:'text'}]},
  action_update: {label:'Atualizar Registro', color:'#6366f1', category:'action', fields:[{key:'modelo',label:'Modelo',type:'text'}]},
  action_delete: {label:'Deletar Registro', color:'#6366f1', category:'action', fields:[{key:'modelo',label:'Modelo',type:'text'}]},
  action_message: {label:'Mostrar Mensagem', color:'#6366f1', category:'action', fields:[{key:'mensagem',label:'Mensagem',type:'text'}]},
  integ_whatsapp: {label:'Enviar WhatsApp', color:'#22c55e', category:'integration', fields:[{key:'para',label:'Para (campo)',type:'text'},{key:'mensagem',label:'Mensagem',type:'text'}]},
  integ_email: {label:'Enviar Email', color:'#f59e0b', category:'integration', fields:[{key:'para',label:'Para (campo)',type:'text'},{key:'assunto',label:'Assunto',type:'text'},{key:'corpo',label:'Corpo',type:'text'}]},
  integ_http: {label:'Chamar API', color:'#06b6d4', category:'integration', fields:[{key:'url',label:'URL',type:'text'},{key:'metodo',label:'Metodo',type:'select',options:['GET','POST','PUT','DELETE']}]},
  logic_if: {label:'Condicao (Se)', color:'#f59e0b', category:'logic', fields:[{key:'campo',label:'Campo',type:'text'},{key:'operador',label:'Operador',type:'select',options:['igual','maior','menor','diferente']},{key:'valor',label:'Valor',type:'text'}]},
  logic_loop: {label:'Para Cada', color:'#f59e0b', category:'logic', fields:[{key:'variavel',label:'Variavel',type:'text'},{key:'colecao',label:'Colecao',type:'text'}]},
  logic_set: {label:'Definir Variavel', color:'#f59e0b', category:'logic', fields:[{key:'nome',label:'Nome',type:'text'},{key:'valor',label:'Valor',type:'text'}]},
  logic_function: {label:'Funcao', color:'#f59e0b', category:'logic', fields:[{key:'nome',label:'Nome',type:'text'},{key:'params',label:'Parametros',type:'text'}]}
};

function dragNode(e, type) {
  e.dataTransfer.setData('nodeType', type);
}

function dropNode(e) {
  e.preventDefault();
  var type = e.dataTransfer.getData('nodeType');
  if (!type || !nodeTypes[type]) return;

  var rect = document.getElementById('flow-canvas').getBoundingClientRect();
  var x = e.clientX - rect.left - 80;
  var y = e.clientY - rect.top - 30;

  var node = {id: ++flowNodeId, type: type, x: x, y: y, props: {}};
  nodeTypes[type].fields.forEach(function(f) { node.props[f.key] = ''; });
  flowNodes.push(node);
  renderFlowNodes();
}

function renderFlowNodes() {
  var container = document.getElementById('flow-nodes');
  container.innerHTML = '';

  flowNodes.forEach(function(node) {
    var nt = nodeTypes[node.type];
    var div = document.createElement('div');
    div.className = 'flow-node' + (node === selectedNode ? ' selected' : '');
    div.style.left = node.x + 'px';
    div.style.top = node.y + 'px';
    div.setAttribute('data-id', node.id);

    var bodyText = '';
    if (nt.fields.length > 0) {
      bodyText = node.props[nt.fields[0].key] || nt.label;
    } else {
      bodyText = nt.label;
    }

    div.innerHTML = '<div class="node-head" style="background:'+nt.color+'15;color:'+nt.color+'">'+
      '<span>'+nt.label+'</span>'+
      '<button onclick="event.stopPropagation();removeNode('+node.id+')" style="background:none;border:none;color:'+nt.color+';cursor:pointer;font-size:14px;opacity:0.5">&times;</button>'+
      '</div>'+
      '<div class="node-body">'+bodyText+'</div>'+
      '<div class="flow-port port-in" onmousedown="event.stopPropagation();startConnect('+node.id+',\'in\',event)" data-node="'+node.id+'" data-port="in"></div>'+
      '<div class="flow-port port-out" onmousedown="event.stopPropagation();startConnect('+node.id+',\'out\',event)" data-node="'+node.id+'" data-port="out"></div>';

    div.onmousedown = function(e) {
      if (e.target.classList.contains('flow-port')) return;
      selectedNode = node;
      draggingNode = node;
      var r = div.getBoundingClientRect();
      dragOffset = {x: e.clientX - r.left, y: e.clientY - r.top};
      renderFlowNodes();
      renderFlowProps();
    };

    container.appendChild(div);
  });

  renderFlowConnections();
}

function flowCanvasMouseDown(e) {
  if (e.target.id === 'flow-canvas' || (e.target.closest('#flow-canvas') && !e.target.closest('.flow-node'))) {
    selectedNode = null;
    renderFlowNodes();
    renderFlowProps();
  }
}

function flowCanvasMouseMove(e) {
  if (draggingNode) {
    var rect = document.getElementById('flow-canvas').getBoundingClientRect();
    draggingNode.x = e.clientX - rect.left - dragOffset.x;
    draggingNode.y = e.clientY - rect.top - dragOffset.y;
    var el = document.querySelector('[data-id="'+draggingNode.id+'"]');
    if (el) {
      el.style.left = draggingNode.x + 'px';
      el.style.top = draggingNode.y + 'px';
    }
    renderFlowConnections();
  }
}

function flowCanvasMouseUp(e) {
  if (connecting && connectFrom) {
    var port = e.target.closest('.flow-port');
    if (port) {
      var toId = parseInt(port.getAttribute('data-node'));
      var toPort = port.getAttribute('data-port');
      if (connectFrom.id !== toId && connectFrom.port !== toPort) {
        var fromId = connectFrom.port === 'out' ? connectFrom.id : toId;
        var targetId = connectFrom.port === 'out' ? toId : connectFrom.id;
        var exists = flowConnections.some(function(c){return c.from===fromId && c.to===targetId;});
        if (!exists) {
          flowConnections.push({from: fromId, to: targetId});
          renderFlowConnections();
          updateFlowCode();
        }
      }
    }
    connecting = false;
    connectFrom = null;
  }
  draggingNode = null;
}

function startConnect(nodeId, port, e) {
  connecting = true;
  connectFrom = {id: nodeId, port: port};
}

function renderFlowConnections() {
  var svg = document.getElementById('flow-svg');
  svg.innerHTML = '';

  flowConnections.forEach(function(conn) {
    var fromEl = document.querySelector('[data-id="'+conn.from+'"] .port-out');
    var toEl = document.querySelector('[data-id="'+conn.to+'"] .port-in');
    if (!fromEl || !toEl) return;

    var canvas = document.getElementById('flow-canvas').getBoundingClientRect();
    var fromRect = fromEl.getBoundingClientRect();
    var toRect = toEl.getBoundingClientRect();

    var x1 = fromRect.left - canvas.left + 6;
    var y1 = fromRect.top - canvas.top + 6;
    var x2 = toRect.left - canvas.left + 6;
    var y2 = toRect.top - canvas.top + 6;

    var dx = Math.abs(x2 - x1) * 0.5;
    var path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
    path.setAttribute('d', 'M'+x1+','+y1+' C'+(x1+dx)+','+y1+' '+(x2-dx)+','+y2+' '+x2+','+y2);
    path.setAttribute('stroke', '#6366f1');
    path.setAttribute('stroke-width', '2');
    path.setAttribute('fill', 'none');
    path.setAttribute('stroke-dasharray', '');
    path.style.pointerEvents = 'stroke';
    path.onclick = function() {
      flowConnections = flowConnections.filter(function(c){return c !== conn;});
      renderFlowConnections();
      updateFlowCode();
    };
    path.style.cursor = 'pointer';
    path.onmouseenter = function(){this.setAttribute('stroke','#ef4444');this.setAttribute('stroke-width','3');};
    path.onmouseleave = function(){this.setAttribute('stroke','#6366f1');this.setAttribute('stroke-width','2');};
    svg.appendChild(path);
  });
}

function removeNode(id) {
  flowNodes = flowNodes.filter(function(n){return n.id !== id;});
  flowConnections = flowConnections.filter(function(c){return c.from !== id && c.to !== id;});
  if (selectedNode && selectedNode.id === id) selectedNode = null;
  renderFlowNodes();
  renderFlowProps();
  updateFlowCode();
}

function renderFlowProps() {
  var panel = document.getElementById('flow-props');
  if (!selectedNode) {
    panel.innerHTML = '<p class="text-xs text-gray-600">Selecione um nodo</p>';
    return;
  }
  var nt = nodeTypes[selectedNode.type];
  var html = '<div class="space-y-3">';
  html += '<div class="text-xs font-semibold mb-2" style="color:'+nt.color+'">'+nt.label+'</div>';

  nt.fields.forEach(function(f) {
    html += '<div><label class="block text-xs text-gray-400 mb-1">'+f.label+'</label>';
    if (f.type === 'select') {
      html += '<select onchange="updateNodeProp(\''+f.key+'\',this.value)" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 focus:outline-none focus:border-primary">';
      f.options.forEach(function(o) {
        html += '<option value="'+o+'"'+(selectedNode.props[f.key]===o?' selected':'')+'>'+o+'</option>';
      });
      html += '</select>';
    } else {
      html += '<input type="text" value="'+(selectedNode.props[f.key]||'')+'" onchange="updateNodeProp(\''+f.key+'\',this.value)" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 focus:outline-none focus:border-primary">';
    }
    html += '</div>';
  });
  html += '</div>';
  panel.innerHTML = html;
}

function updateNodeProp(key, val) {
  if (!selectedNode) return;
  selectedNode.props[key] = val;
  renderFlowNodes();
  updateFlowCode();
}

function updateFlowCode() {
  var code = '';
  var triggers = flowNodes.filter(function(n){return nodeTypes[n.type].category==='trigger';});

  if (triggers.length > 0) {
    code += 'eventos\n\n';
    triggers.forEach(function(trigger) {
      switch(trigger.type) {
        case 'trigger_click':
          code += '  quando clicar "' + (trigger.props.botao||'Botao') + '"\n';
          break;
        case 'trigger_create':
          code += '  quando criar ' + (trigger.props.modelo||'item') + '\n';
          break;
        case 'trigger_update':
          code += '  quando atualizar ' + (trigger.props.modelo||'item') + '\n';
          break;
        case 'trigger_cron':
          code += '  cada ' + (trigger.props.intervalo||'5') + ' ' + (trigger.props.unidade||'minutos') + '\n';
          break;
      }
      var connected = getConnectedNodes(trigger.id);
      connected.forEach(function(node) {
        code += '    ' + nodeToCode(node) + '\n';
      });
      code += '\n';
    });
  }

  var logicNodes = flowNodes.filter(function(n){
    return nodeTypes[n.type].category==='logic' && !flowConnections.some(function(c){return c.to===n.id;});
  });
  if (logicNodes.length > 0) {
    code += 'logica\n\n';
    logicNodes.forEach(function(n) {
      code += '  ' + nodeToCode(n) + '\n';
      var connected = getConnectedNodes(n.id);
      connected.forEach(function(cn) {
        code += '    ' + nodeToCode(cn) + '\n';
      });
    });
  }

  document.getElementById('flow-generated-code').textContent = code;
}

function getConnectedNodes(fromId) {
  var result = [];
  flowConnections.filter(function(c){return c.from===fromId;}).forEach(function(c) {
    var node = flowNodes.find(function(n){return n.id===c.to;});
    if (node) result.push(node);
  });
  return result;
}

function nodeToCode(node) {
  switch(node.type) {
    case 'action_create': return 'criar ' + (node.props.modelo||'item');
    case 'action_update': return 'atualizar ' + (node.props.modelo||'item');
    case 'action_delete': return 'deletar ' + (node.props.modelo||'item');
    case 'action_message': return 'mostrar "' + (node.props.mensagem||'') + '"';
    case 'integ_whatsapp': return 'enviar mensagem para ' + (node.props.para||'telefone') + ' texto "' + (node.props.mensagem||'') + '"';
    case 'integ_email': return 'enviar email para ' + (node.props.para||'email') + ' assunto "' + (node.props.assunto||'') + '"';
    case 'integ_http': return 'chamar api "' + (node.props.url||'') + '"';
    case 'logic_if': return 'se ' + (node.props.campo||'x') + ' ' + (node.props.operador||'igual') + ' ' + (node.props.valor||'0');
    case 'logic_loop': return 'para cada ' + (node.props.variavel||'item') + ' em ' + (node.props.colecao||'lista');
    case 'logic_set': return 'definir ' + (node.props.nome||'x') + ' = ' + (node.props.valor||'0');
    case 'logic_function': return 'funcao ' + (node.props.nome||'minha_funcao') + '(' + (node.props.params||'') + ')';
    default: return '# ' + node.type;
  }
}

function generateFromFlow() {
  updateFlowCode();
  var code = document.getElementById('flow-generated-code').textContent;
  if (!code) { termLog('error', 'Nenhum fluxo criado'); return; }

  fetch('/api/file/save', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({path: 'fluxo_visual.fg', content: code})
  }).then(function() {
    termLog('success', 'Fluxo salvo em fluxo_visual.fg');
    loadFileTree();
    openFile('fluxo_visual.fg');
    switchMode('editor');
  });
}

function clearFlow() {
  flowNodes = [];
  flowConnections = [];
  selectedNode = null;
  renderFlowNodes();
  renderFlowProps();
  document.getElementById('flow-generated-code').textContent = '';
}
</script>
</body>
</html>`
