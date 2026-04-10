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

  <!-- Editor area -->
  <div class="flex-1 flex flex-col overflow-hidden">
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
</script>
</body>
</html>`
