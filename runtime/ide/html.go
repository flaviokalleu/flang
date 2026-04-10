package ide

var ideHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Flang IDE</title>
<script src="https://cdn.tailwindcss.com"></script>
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

  <!-- Panel: Designer (visual screen builder) -->
  <div id="panel-designer" style="display:none" class="flex flex-1 overflow-hidden">

    <!-- Component Palette -->
    <div class="w-48 bg-gray-900 border-r border-gray-800 p-3 overflow-y-auto flex-shrink-0">
      <div class="text-xs font-semibold text-gray-500 uppercase mb-3">Componentes</div>
      <div class="space-y-1">
        <div class="palette-item" draggable="true" ondragstart="dragComp(event,'titulo')" data-type="titulo">
          <span class="w-4 h-4">T</span> Titulo
        </div>
        <div class="palette-item" draggable="true" ondragstart="dragComp(event,'lista')" data-type="lista">
          <span class="w-4 h-4">&#9776;</span> Lista/Tabela
        </div>
        <div class="palette-item" draggable="true" ondragstart="dragComp(event,'botao')" data-type="botao">
          <span class="w-4 h-4">&#9634;</span> Botao
        </div>
        <div class="palette-item" draggable="true" ondragstart="dragComp(event,'busca')" data-type="busca">
          <span class="w-4 h-4">&#128269;</span> Busca
        </div>
        <div class="palette-item" draggable="true" ondragstart="dragComp(event,'grafico')" data-type="grafico">
          <span class="w-4 h-4">&#128202;</span> Grafico
        </div>
        <div class="palette-item" draggable="true" ondragstart="dragComp(event,'dashboard')" data-type="dashboard">
          <span class="w-4 h-4">&#128203;</span> Dashboard
        </div>
        <div class="palette-item" draggable="true" ondragstart="dragComp(event,'formulario')" data-type="formulario">
          <span class="w-4 h-4">&#128221;</span> Formulario
        </div>
        <div class="palette-item" draggable="true" ondragstart="dragComp(event,'texto')" data-type="texto">
          <span class="w-4 h-4">Aa</span> Texto
        </div>
      </div>

      <div class="text-xs font-semibold text-gray-500 uppercase mt-6 mb-3">Telas</div>
      <button onclick="addScreen()" class="w-full px-3 py-2 text-xs bg-primary/10 text-primary rounded-lg hover:bg-primary/20 transition-all">+ Nova Tela</button>
      <div id="screen-list" class="mt-2 space-y-1"></div>
    </div>

    <!-- Canvas -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <div class="px-4 py-2 bg-gray-900 border-b border-gray-800 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="text-sm font-semibold" id="canvas-screen-name">Selecione uma tela</span>
        </div>
        <button onclick="generateFromDesigner()" class="px-3 py-1.5 text-xs bg-primary hover:bg-primary/80 text-white rounded-lg transition-all">Gerar Codigo .fg</button>
      </div>
      <div id="canvas" class="flex-1 overflow-y-auto p-6 bg-gray-950"
           ondragover="event.preventDefault();this.classList.add('ring-2','ring-primary/50')"
           ondragleave="this.classList.remove('ring-2','ring-primary/50')"
           ondrop="dropComp(event);this.classList.remove('ring-2','ring-primary/50')">
        <div id="canvas-empty" class="flex flex-col items-center justify-center h-full text-gray-600">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" class="w-16 h-16 mb-4 opacity-30"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>
          <p class="text-sm">Arraste componentes aqui</p>
          <p class="text-xs mt-1 opacity-50">ou crie uma nova tela na esquerda</p>
        </div>
      </div>
    </div>

    <!-- Properties Panel -->
    <div class="w-64 bg-gray-900 border-l border-gray-800 p-4 overflow-y-auto flex-shrink-0">
      <div class="text-xs font-semibold text-gray-500 uppercase mb-3">Propriedades</div>
      <div id="props-panel">
        <p class="text-xs text-gray-600">Selecione um componente</p>
      </div>

      <div class="text-xs font-semibold text-gray-500 uppercase mt-6 mb-3">Codigo Gerado</div>
      <pre id="generated-code" class="text-xs bg-gray-950 text-green-400 p-3 rounded-lg overflow-auto max-h-60 font-mono"></pre>
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
// DESIGNER - Visual Screen Builder
// ============================================================

var screens = [];
var currentScreen = null;
var selectedComp = null;
var compIdCounter = 0;

function addScreen() {
  var name = prompt('Nome da tela:');
  if (!name) return;
  var screen = {name: name, title: name.charAt(0).toUpperCase()+name.slice(1), components: []};
  screens.push(screen);
  currentScreen = screen;
  renderScreenList();
  renderCanvas();
  document.getElementById('canvas-screen-name').textContent = screen.title;
}

function selectScreen(idx) {
  currentScreen = screens[idx];
  renderScreenList();
  renderCanvas();
  document.getElementById('canvas-screen-name').textContent = currentScreen.title;
}

function renderScreenList() {
  var html = '';
  screens.forEach(function(s, i) {
    var active = s === currentScreen ? 'bg-primary/10 text-primary' : 'text-gray-400 hover:bg-gray-800';
    html += '<div class="px-3 py-1.5 text-xs rounded-lg cursor-pointer '+active+'" onclick="selectScreen('+i+')">'+s.title+'</div>';
  });
  document.getElementById('screen-list').innerHTML = html;
}

function dragComp(e, type) {
  e.dataTransfer.setData('compType', type);
}

function dropComp(e) {
  e.preventDefault();
  var type = e.dataTransfer.getData('compType');
  if (!type || !currentScreen) return;

  var comp = {id: ++compIdCounter, type: type, props: {}};

  switch(type) {
    case 'titulo': comp.props = {texto: 'Minha Tela'}; break;
    case 'lista': comp.props = {modelo: 'produto', campos: 'nome, preco, status'}; break;
    case 'botao': comp.props = {texto: 'Novo', cor: 'azul', acao: 'criar produto'}; break;
    case 'busca': comp.props = {modelo: 'produto'}; break;
    case 'grafico': comp.props = {modelo: 'produto', tipo: 'barra'}; break;
    case 'dashboard': comp.props = {}; break;
    case 'formulario': comp.props = {modelo: 'produto'}; break;
    case 'texto': comp.props = {conteudo: 'Texto aqui...'}; break;
  }

  currentScreen.components.push(comp);
  renderCanvas();
  updateGeneratedCode();
}

function renderCanvas() {
  var canvas = document.getElementById('canvas');
  var empty = document.getElementById('canvas-empty');

  if (!currentScreen || !currentScreen.components.length) {
    empty.style.display = 'flex';
    Array.from(canvas.children).forEach(function(c) {
      if (c.id !== 'canvas-empty') c.remove();
    });
    return;
  }
  empty.style.display = 'none';

  Array.from(canvas.children).forEach(function(c) {
    if (c.id !== 'canvas-empty') c.remove();
  });

  currentScreen.components.forEach(function(comp, idx) {
    var div = document.createElement('div');
    div.className = 'canvas-comp' + (comp === selectedComp ? ' selected' : '');
    div.onclick = function(e) { e.stopPropagation(); selectComp(idx); };

    var preview = '';
    switch(comp.type) {
      case 'titulo':
        preview = '<div style="font-size:18px;font-weight:700">' + (comp.props.texto||'Titulo') + '</div>';
        break;
      case 'lista':
        preview = '<div style="display:flex;gap:8px;opacity:0.5"><span>ID</span><span>|</span>' +
          (comp.props.campos||'nome').split(',').map(function(c){return '<span>'+c.trim()+'</span>';}).join('<span>|</span>') +
          '</div><div style="border-top:1px solid rgba(255,255,255,0.05);margin-top:6px;padding-top:6px;opacity:0.3;font-size:11px">dados carregam automaticamente</div>';
        break;
      case 'botao':
        preview = '<div style="display:inline-block;padding:6px 16px;background:rgba(99,102,241,0.2);border-radius:8px;font-size:12px;color:#818cf8">' + (comp.props.texto||'Botao') + '</div>';
        break;
      case 'busca':
        preview = '<div style="padding:6px 12px;background:rgba(255,255,255,0.03);border:1px solid rgba(255,255,255,0.06);border-radius:8px;font-size:12px;opacity:0.5">Buscar em ' + (comp.props.modelo||'...') + '</div>';
        break;
      case 'grafico':
        preview = '<div style="height:40px;display:flex;align-items:flex-end;gap:3px"><div style="width:16px;background:rgba(99,102,241,0.3);height:60%;border-radius:3px 3px 0 0"></div><div style="width:16px;background:rgba(99,102,241,0.5);height:80%;border-radius:3px 3px 0 0"></div><div style="width:16px;background:rgba(99,102,241,0.4);height:50%;border-radius:3px 3px 0 0"></div><div style="width:16px;background:rgba(99,102,241,0.6);height:90%;border-radius:3px 3px 0 0"></div></div>';
        break;
      case 'dashboard':
        preview = '<div style="display:grid;grid-template-columns:1fr 1fr;gap:6px"><div style="padding:8px;background:rgba(99,102,241,0.1);border-radius:6px;font-size:10px">Card 1</div><div style="padding:8px;background:rgba(34,197,94,0.1);border-radius:6px;font-size:10px">Card 2</div></div>';
        break;
      case 'formulario':
        preview = '<div style="space-y:4px"><div style="height:8px;width:60%;background:rgba(255,255,255,0.05);border-radius:4px;margin-bottom:4px"></div><div style="height:28px;background:rgba(255,255,255,0.03);border:1px solid rgba(255,255,255,0.06);border-radius:6px;margin-bottom:4px"></div><div style="height:8px;width:40%;background:rgba(255,255,255,0.05);border-radius:4px;margin-bottom:4px"></div><div style="height:28px;background:rgba(255,255,255,0.03);border:1px solid rgba(255,255,255,0.06);border-radius:6px"></div></div>';
        break;
      case 'texto':
        preview = '<div style="font-size:13px;opacity:0.6">' + (comp.props.conteudo||'Texto') + '</div>';
        break;
    }

    div.innerHTML = '<div class="comp-label">' + comp.type + '</div>' +
      '<div class="comp-preview">' + preview + '</div>' +
      '<button class="comp-delete" onclick="event.stopPropagation();removeComp('+idx+')">&times;</button>';
    canvas.appendChild(div);
  });
}

function selectComp(idx) {
  selectedComp = currentScreen.components[idx];
  renderCanvas();
  renderProps();
}

function removeComp(idx) {
  currentScreen.components.splice(idx, 1);
  selectedComp = null;
  renderCanvas();
  renderProps();
  updateGeneratedCode();
}

function renderProps() {
  var panel = document.getElementById('props-panel');
  if (!selectedComp) {
    panel.innerHTML = '<p class="text-xs text-gray-600">Selecione um componente</p>';
    return;
  }

  var html = '<div class="space-y-3">';
  html += '<div class="text-xs text-primary font-semibold mb-2">' + selectedComp.type.toUpperCase() + '</div>';

  var fields = {
    titulo: [{key:'texto',label:'Texto',type:'text'}],
    lista: [{key:'modelo',label:'Modelo',type:'text'},{key:'campos',label:'Campos (virgula)',type:'text'}],
    botao: [{key:'texto',label:'Texto',type:'text'},{key:'cor',label:'Cor',type:'select',options:['azul','verde','vermelho','amarelo']},{key:'acao',label:'Acao',type:'text'}],
    busca: [{key:'modelo',label:'Modelo',type:'text'}],
    grafico: [{key:'modelo',label:'Modelo',type:'text'},{key:'tipo',label:'Tipo',type:'select',options:['barra','pizza','linha']}],
    formulario: [{key:'modelo',label:'Modelo',type:'text'}],
    texto: [{key:'conteudo',label:'Conteudo',type:'textarea'}],
    dashboard: []
  };

  (fields[selectedComp.type]||[]).forEach(function(f) {
    html += '<div><label class="block text-xs text-gray-400 mb-1">' + f.label + '</label>';
    if (f.type === 'select') {
      html += '<select onchange="updateProp(\''+f.key+'\',this.value)" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 focus:outline-none focus:border-primary">';
      f.options.forEach(function(o) {
        html += '<option value="'+o+'"'+(selectedComp.props[f.key]===o?' selected':'')+'>'+o+'</option>';
      });
      html += '</select>';
    } else if (f.type === 'textarea') {
      html += '<textarea onchange="updateProp(\''+f.key+'\',this.value)" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 focus:outline-none focus:border-primary resize-none" rows="3">'+(selectedComp.props[f.key]||'')+'</textarea>';
    } else {
      html += '<input type="text" value="'+(selectedComp.props[f.key]||'')+'" onchange="updateProp(\''+f.key+'\',this.value)" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-1.5 text-xs text-gray-200 focus:outline-none focus:border-primary">';
    }
    html += '</div>';
  });

  html += '</div>';
  panel.innerHTML = html;
}

function updateProp(key, val) {
  if (!selectedComp) return;
  selectedComp.props[key] = val;
  renderCanvas();
  updateGeneratedCode();
}

function updateGeneratedCode() {
  var code = '';
  screens.forEach(function(s) {
    code += 'telas\n\n';
    code += '  tela ' + s.name + '\n';
    s.components.forEach(function(c) {
      switch(c.type) {
        case 'titulo': code += '    titulo "' + (c.props.texto||'') + '"\n'; break;
        case 'lista':
          code += '    lista ' + (c.props.modelo||'item') + '\n';
          (c.props.campos||'').split(',').forEach(function(f) {
            f = f.trim();
            if (f) code += '      mostrar ' + f + '\n';
          });
          break;
        case 'botao': code += '    botao ' + (c.props.cor||'azul') + '\n      texto "' + (c.props.texto||'Novo') + '"\n'; break;
        case 'busca': code += '    busca ' + (c.props.modelo||'item') + '\n'; break;
        case 'grafico': code += '    grafico ' + (c.props.modelo||'item') + '\n      tipo ' + (c.props.tipo||'barra') + '\n'; break;
        case 'dashboard': code += '    dashboard\n'; break;
        case 'formulario': code += '    formulario ' + (c.props.modelo||'item') + '\n'; break;
        case 'texto': code += '    texto "' + (c.props.conteudo||'') + '"\n'; break;
      }
    });
    code += '\n';
  });
  document.getElementById('generated-code').textContent = code;
}

function generateFromDesigner() {
  updateGeneratedCode();
  var code = document.getElementById('generated-code').textContent;
  if (!code) { termLog('error', 'Nenhuma tela criada'); return; }

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
