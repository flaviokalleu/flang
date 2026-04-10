package servidor

import (
	"fmt"
	"html"
	"strings"

	"github.com/flavio/flang/compiler/ast"
)

// renderHTML generates the full single-page application HTML.
// It uses CSS variables derived from the Theme so users can control
// every visual aspect via their .fg file.
func (s *Servidor) renderHTML() string {
	theme := s.Program.Theme
	if theme == nil {
		theme = ast.DefaultTheme()
	}
	applyThemeDefaults(theme)

	darkClass := ""
	if theme.Dark {
		darkClass = " dark"
	}

	var b strings.Builder

	// --- Head ---
	b.WriteString(`<!DOCTYPE html><html lang="pt-BR"><head><meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>` + html.EscapeString(cap(s.Program.System.Name)) + `</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=` + strings.ReplaceAll(theme.Font, " ", "+") + `:wght@300;400;500;600;700;800&display=swap" rel="stylesheet">
<script src="https://cdn.jsdelivr.net/npm/chart.js@4/dist/chart.umd.min.js"></script>
<style>`)
	b.WriteString(s.generateCSS(theme))
	b.WriteString(`</style></head><body class="` + darkClass + `">`)

	// --- Sidebar ---
	b.WriteString(s.renderSidebar(theme))

	// --- Main area ---
	b.WriteString(`<div class="main" id="main">`)
	b.WriteString(s.renderTopbar())
	b.WriteString(`<div class="content">`)

	// Dashboard section
	b.WriteString(s.renderDashboard(theme))

	// Model sections - respect custom screens
	if len(s.Program.Screens) > 0 {
		s.renderCustomScreens(&b)
	}
	// Always generate model sections (custom screens are shown via nav)
	for _, model := range s.Program.Models {
		s.renderModelSection(&b, model, theme)
	}

	b.WriteString(`</div></div>`) // content, main

	// Toast container
	b.WriteString(`<div id="toast" class="toast"></div>`)

	// Auth modal (login/register)
	if s.Auth != nil {
		loginField := "email"
		passField := "senha"
		if s.Program.Auth != nil {
			if s.Program.Auth.LoginField != "" {
				loginField = s.Program.Auth.LoginField
			}
			if s.Program.Auth.PassField != "" {
				passField = s.Program.Auth.PassField
			}
		}
		b.WriteString(`<div id="auth-modal" class="modal-overlay" style="display:none">`)
		b.WriteString(`<div class="modal card" style="max-width:400px;padding:32px">`)
		b.WriteString(`<h2 id="auth-title" style="margin:0 0 20px;text-align:center">Entrar</h2>`)
		b.WriteString(`<form id="auth-form" onsubmit="authSubmit(event)">`)
		b.WriteString(`<div id="auth-extra-fields"></div>`)
		b.WriteString(fmt.Sprintf(`<div class="form-group"><label>%s</label><input type="text" id="auth-login" required class="form-input" placeholder="%s"></div>`, cap(loginField), loginField))
		b.WriteString(fmt.Sprintf(`<div class="form-group"><label>%s</label><input type="password" id="auth-pass" required class="form-input" placeholder="%s" minlength="6"></div>`, cap(passField), passField))
		b.WriteString(`<div id="auth-error" style="color:#ef4444;font-size:13px;margin:8px 0;display:none"></div>`)
		b.WriteString(`<button type="submit" class="btn btn-glow" style="width:100%;margin-top:12px">Entrar</button>`)
		b.WriteString(`<p style="text-align:center;margin:16px 0 0;font-size:13px;opacity:.7">`)
		b.WriteString(`<span id="auth-toggle-text">Nao tem conta?</span> <a href="#" onclick="toggleAuthMode()" id="auth-toggle-link" style="color:var(--primary)">Criar conta</a></p>`)
		b.WriteString(`<button type="button" onclick="fecharAuth()" class="btn" style="width:100%;margin-top:8px;background:transparent;border:1px solid rgba(255,255,255,.1)">Cancelar</button>`)
		b.WriteString(`</form></div></div>`)
	}

	// --- JavaScript ---
	b.WriteString(`<script>`)
	b.WriteString(s.generateJS(theme))
	b.WriteString(`</script></body></html>`)

	return b.String()
}

// applyThemeDefaults fills empty theme fields with sensible defaults.
func applyThemeDefaults(t *ast.Theme) {
	if t.Primary == "" {
		t.Primary = "#6366f1"
	}
	if t.Secondary == "" {
		t.Secondary = "#8b5cf6"
	}
	if t.Accent == "" {
		t.Accent = "#f59e0b"
	}
	if t.Sidebar == "" {
		t.Sidebar = "#1e1b4b"
	}
	if t.Font == "" {
		t.Font = "Inter"
	}
	if t.Radius == "" {
		t.Radius = "12px"
	}
	if t.Style == "" {
		t.Style = "glassmorphism"
	}
	if t.Background == "" {
		if t.Dark {
			t.Background = "#0c0a1d"
		} else {
			t.Background = "#f8fafc"
		}
	}
	if t.CardBg == "" {
		if t.Dark {
			t.CardBg = "rgba(30,27,75,0.6)"
		} else {
			t.CardBg = "rgba(255,255,255,0.85)"
		}
	}
	if t.TextColor == "" {
		if t.Dark {
			t.TextColor = "#e2e8f0"
		} else {
			t.TextColor = "#0f172a"
		}
	}
}

// ============================================================
// Sidebar
// ============================================================

func (s *Servidor) renderSidebar(theme *ast.Theme) string {
	var b strings.Builder
	b.WriteString(`<aside class="sidebar" id="sidebar">`)
	b.WriteString(`<div class="sb-top">`)
	// Brand
	b.WriteString(`<div class="sb-brand">`)
	if theme.Icon != "" {
		b.WriteString(`<div class="sb-logo"><img src="` + theme.Icon + `" alt="logo" style="width:24px;height:24px;object-fit:contain"></div>`)
	} else {
		b.WriteString(`<div class="sb-logo">` + svgIcon("zap") + `</div>`)
	}
	b.WriteString(`<span class="sb-name">` + html.EscapeString(cap(s.Program.System.Name)) + `</span>`)
	b.WriteString(`<button class="sb-collapse" onclick="toggleCollapse()" title="Recolher">` + svgIcon("chevleft") + `</button></div>`)
	// Nav
	b.WriteString(`<nav class="sb-nav">`)
	if len(s.Program.SidebarItems) > 0 {
		// Custom sidebar items defined by user
		for _, item := range s.Program.SidebarItems {
			label := item.Label
			if label == "" {
				continue
			}
			icon := item.Icon
			if icon == "" {
				icon = "grid"
			}
			link := item.Link
			if link == "" {
				link = lo(label)
			}
			b.WriteString(fmt.Sprintf(`<a class="sb-link" onclick="irPara('%s',this)" href="#">`, html.EscapeString(link)))
			b.WriteString(`<div class="sb-icon">` + svgIcon(icon) + `</div><span>` + html.EscapeString(cap(label)) + `</span></a>`)
		}
	} else {
		// Smart sidebar: if user defined screens, use those; otherwise auto-generate from models
		b.WriteString(`<a class="sb-link active" onclick="irPara('dashboard',this)" href="#">`)
		b.WriteString(`<div class="sb-icon">` + svgIcon("layout") + `</div><span>Dashboard</span></a>`)

		if len(s.Program.Screens) > 0 {
			// User defined screens — show those as primary navigation
			for _, scr := range s.Program.Screens {
				name := lo(scr.Name)
				title := scr.Title
				if title == "" {
					title = cap(scr.Name)
				}
				icon := "grid"
				// Try to infer icon from the model referenced by the screen
				if pm := s.primaryScreenModel(scr); pm != nil {
					icon = modelIcon(lo(pm.Name))
				}
				b.WriteString(fmt.Sprintf(`<a class="sb-link" onclick="irPara('screen-%s',this)" href="#">`, name))
				b.WriteString(`<div class="sb-icon">` + svgIcon(icon) + `</div><span>` + html.EscapeString(title) + `</span></a>`)
			}
			// Add models that don't have a custom screen
			for _, model := range s.Program.Models {
				mName := lo(model.Name)
				hasScreen := false
				for _, scr := range s.Program.Screens {
					if pm := s.primaryScreenModel(scr); pm != nil && lo(pm.Name) == mName {
						hasScreen = true
						break
					}
				}
				if !hasScreen {
					icon := modelIcon(mName)
					b.WriteString(fmt.Sprintf(`<a class="sb-link" onclick="irPara('%s',this)" href="#">`, mName))
					b.WriteString(`<div class="sb-icon">` + svgIcon(icon) + `</div><span>` + html.EscapeString(cap(model.Name)) + `</span></a>`)
				}
			}
		} else {
			// No custom screens — auto-generate from models
			for _, model := range s.Program.Models {
				name := lo(model.Name)
				icon := modelIcon(name)
				if model.Icon != "" {
					icon = model.Icon
				}
				b.WriteString(fmt.Sprintf(`<a class="sb-link" onclick="irPara('%s',this)" href="#">`, name))
				b.WriteString(`<div class="sb-icon">` + svgIcon(icon) + `</div><span>` + html.EscapeString(cap(model.Name)) + `</span></a>`)
			}
		}
	}
	b.WriteString(`</nav></div>`)
	// Footer
	b.WriteString(`<div class="sb-foot">`)
	b.WriteString(`<button class="sb-theme" onclick="toggleDark()">` + svgIcon("moon") + `<span>Tema</span></button>`)
	b.WriteString(`<div class="sb-powered">Flang v0.3</div>`)
	b.WriteString(`</div></aside>`)
	return b.String()
}

// ============================================================
// Topbar
// ============================================================

func (s *Servidor) renderTopbar() string {
	var b strings.Builder
	b.WriteString(`<header class="topbar">`)
	b.WriteString(`<button class="tb-menu" onclick="toggleSidebar()">` + svgIcon("menu") + `</button>`)
	b.WriteString(`<div class="tb-title-wrap"><h1 id="page-title">Dashboard</h1><p class="tb-sub">Operacao em tempo real com Flang</p></div>`)
	b.WriteString(`<div class="tb-status">`)
	b.WriteString(`<div class="tb-chip"><span class="tb-dot tb-dot-online"></span><span id="tb-sockets">0 conexoes</span></div>`)
	b.WriteString(`<div class="tb-chip"><span class="tb-dot tb-dot-jobs"></span><span id="tb-jobs">jobs 0/0</span></div>`)
	b.WriteString(`<div class="tb-chip"><span class="tb-dot tb-dot-wa"></span><span id="tb-wa">whatsapp offline</span></div>`)
	b.WriteString(`</div>`)
	b.WriteString(`<div class="tb-end">`)
	b.WriteString(`<div class="tb-search"><input type="text" placeholder="Buscar..." id="global-search" oninput="buscaGlobal(this.value)">` + svgIcon("search") + `</div>`)
	// Auth button
	if s.Auth != nil {
		b.WriteString(`<div id="auth-area">`)
		b.WriteString(`<button class="btn btn-sm" id="btn-login" onclick="mostrarLogin()" style="margin-left:8px">` + svgIcon("user") + ` <span>Entrar</span></button>`)
		b.WriteString(`<span id="user-info" style="display:none;margin-left:8px;font-size:13px;opacity:.8"></span>`)
		b.WriteString(`<button class="btn btn-sm" id="btn-logout" onclick="sair()" style="display:none;margin-left:4px">Sair</button>`)
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div></header>`)
	return b.String()
}

// ============================================================
// Dashboard
// ============================================================

func (s *Servidor) renderDashboard(theme *ast.Theme) string {
	var b strings.Builder
	b.WriteString(`<div id="secao-dashboard" class="section anim-in">`)

	// Bento stat cards
	b.WriteString(`<div class="bento">`)
	colors := []string{theme.Primary, theme.Secondary, theme.Accent, "#10b981", "#3b82f6", "#ef4444", "#06b6d4", "#ec4899"}
	for i, model := range s.Program.Models {
		name := lo(model.Name)
		icon := modelIcon(name)
		color := colors[i%len(colors)]
		b.WriteString(fmt.Sprintf(`<div class="bento-card" onclick="irParaNav('%s')" style="--accent:%s">`, name, color))
		b.WriteString(`<div class="bc-icon">` + svgIcon(icon) + `</div>`)
		b.WriteString(fmt.Sprintf(`<div class="bc-num" id="stat-%s">0</div>`, name))
		b.WriteString(`<div class="bc-label">` + html.EscapeString(cap(model.Name)) + `</div>`)
		b.WriteString(`<div class="bc-glow"></div></div>`)
	}
	b.WriteString(`</div>`)

	// Chart.js canvas for records per model
	b.WriteString(`<div class="card chart-card"><div class="card-head">` + svgIcon("activity") + `<h3>Registros por Modelo</h3></div>`)
	b.WriteString(`<div class="chart-wrap"><canvas id="chart-models" height="260"></canvas></div></div>`)

	// Render any user-defined grafico components from screens
	for _, scr := range s.Program.Screens {
		for _, comp := range scr.Components {
			if comp.Type == ast.CompChart {
				s.renderChartComponent(&b, comp)
			}
		}
	}

	// Status chart for models with status fields
	hasStatus := false
	for _, model := range s.Program.Models {
		for _, f := range model.Fields {
			if f.Type == ast.FieldStatus {
				hasStatus = true
				break
			}
		}
	}
	if hasStatus {
		b.WriteString(`<div class="card chart-card"><div class="card-head">` + svgIcon("tag") + `<h3>Status por Modelo</h3></div>`)
		b.WriteString(`<div class="chart-wrap"><canvas id="chart-status" height="260"></canvas></div></div>`)
	}

	// Activity feed + info
	b.WriteString(`<div class="dash-grid">`)
	b.WriteString(`<div class="card"><div class="card-head">` + svgIcon("activity") + `<h3>Atividade Recente</h3></div>`)
	b.WriteString(`<div id="atividade" class="ativ-list"><div class="empty-state">` + svgIcon("inbox") + `<p>Nenhuma atividade</p></div></div></div>`)
	b.WriteString(`<div class="card"><div class="card-head">` + svgIcon("info") + `<h3>Informa&ccedil;&otilde;es</h3></div>`)
	b.WriteString(`<div class="info-list">`)
	b.WriteString(fmt.Sprintf(`<div class="info-row"><span class="info-k">Sistema</span><span class="info-v">%s</span></div>`, html.EscapeString(cap(s.Program.System.Name))))
	b.WriteString(fmt.Sprintf(`<div class="info-row"><span class="info-k">Modelos</span><span class="info-v">%d</span></div>`, len(s.Program.Models)))
	b.WriteString(fmt.Sprintf(`<div class="info-row"><span class="info-k">Telas</span><span class="info-v">%d</span></div>`, len(s.Program.Screens)))
	b.WriteString(`<div class="info-row"><span class="info-k">Engine</span><span class="info-v">Flang v0.3</span></div>`)
	b.WriteString(`</div></div>`)
	b.WriteString(`</div>`) // dash-grid
	b.WriteString(`</div>`) // dashboard section
	return b.String()
}

// ============================================================
// Chart component (user-defined grafico blocks)
// ============================================================

func (s *Servidor) renderChartComponent(b *strings.Builder, comp *ast.Component) {
	chartType := comp.Properties["tipo"]
	if chartType == "" {
		chartType = "bar"
	}
	target := comp.Target
	title := comp.Properties["titulo"]
	if title == "" {
		title = cap(target) + " - Grafico"
	}
	chartID := "chart-custom-" + lo(target) + "-" + lo(chartType)
	b.WriteString(`<div class="card chart-card"><div class="card-head">` + svgIcon("activity") + `<h3>` + title + `</h3></div>`)
	b.WriteString(fmt.Sprintf(`<div class="chart-wrap"><canvas id="%s" height="260" data-chart-type="%s" data-chart-model="%s"></canvas></div></div>`, chartID, chartType, lo(target)))
}

// ============================================================
// Custom screens
// ============================================================

func (s *Servidor) renderCustomScreens(b *strings.Builder) {
	for _, scr := range s.Program.Screens {
		name := lo(scr.Name)
		title := scr.Title
		primaryModel := s.primaryScreenModel(scr)
		if title == "" {
			title = cap(scr.Name)
		}
		b.WriteString(fmt.Sprintf(`<div id="secao-screen-%s" class="section" style="display:none">`, name))
		b.WriteString(`<div class="sec-head"><div class="sec-left"><h2>` + html.EscapeString(title) + `</h2></div></div>`)
		for _, comp := range scr.Components {
			s.renderScreenComponent(b, comp, primaryModel)
		}
		b.WriteString(`</div>`)
	}
}

func (s *Servidor) primaryScreenModel(screen *ast.Screen) *ast.Model {
	if screen == nil {
		return nil
	}
	for _, comp := range screen.Components {
		if model := s.inferModelForTarget(comp.Target); model != nil {
			return model
		}
	}
	return nil
}

func (s *Servidor) findModelByName(name string) *ast.Model {
	target := lo(name)
	for _, model := range s.Program.Models {
		if lo(model.Name) == target {
			return model
		}
	}
	return nil
}

func (s *Servidor) inferModelForTarget(target string) *ast.Model {
	if target == "" {
		return nil
	}
	if model := s.findModelByName(target); model != nil {
		return model
	}
	lookup := lo(target)
	if strings.HasSuffix(lookup, "s") {
		return s.findModelByName(strings.TrimSuffix(lookup, "s"))
	}
	return nil
}

func (s *Servidor) renderInlineListComponent(b *strings.Builder, model *ast.Model) {
	if model == nil {
		return
	}

	name := lo(model.Name)
	b.WriteString(fmt.Sprintf(`<div class="card table-wrap" data-list-model="%s">`, name))
	b.WriteString(`<table><thead><tr><th class="th-id">#</th>`)
	for _, f := range model.Fields {
		if f.Type == ast.FieldSenha {
			continue
		}
		b.WriteString(`<th>` + html.EscapeString(cap(f.Name)) + `</th>`)
	}
	b.WriteString(`<th class="th-act"></th></tr></thead><tbody></tbody></table>`)
	b.WriteString(`<div class="empty-state">`)
	b.WriteString(svgIcon("inbox") + `<p>Nenhum registro</p></div></div>`)
}

func (s *Servidor) renderScreenComponent(b *strings.Builder, comp *ast.Component, primaryModel *ast.Model) {
	switch comp.Type {
	case ast.CompList:
		s.renderInlineListComponent(b, s.inferModelForTarget(comp.Target))
	case ast.CompChart:
		s.renderChartComponent(b, comp)
	case ast.CompText:
		text := comp.Properties["conteudo"]
		if text == "" {
			text = comp.Properties["valor"]
		}
		b.WriteString(`<div class="card" style="padding:20px"><p>` + text + `</p></div>`)
	case ast.CompButton:
		label := comp.Properties["texto"]
		if label == "" {
			label = comp.Properties["text"]
		}
		if label == "" {
			label = comp.Properties["label"]
		}
		action := comp.Properties["acao"]
		if action == "" {
			if primaryModel != nil {
				action = fmt.Sprintf("abrirForm('%s')", lo(primaryModel.Name))
			} else if model := s.inferModelForTarget(comp.Target); model != nil {
				action = fmt.Sprintf("abrirForm('%s')", lo(model.Name))
			}
		}
		b.WriteString(fmt.Sprintf(`<button class="btn btn-glow" onclick="%s">%s</button>`, action, html.EscapeString(label)))
	case ast.CompChat:
		s.renderChatComponent(b, comp)
	case ast.CompForm:
		target := lo(comp.Target)
		b.WriteString(fmt.Sprintf(`<div class="card" style="padding:20px"><h3>Formulario - %s</h3>`, cap(target)))
		b.WriteString(fmt.Sprintf(`<form onsubmit="salvar('%s',event)" class="modal-form">`, target))
		b.WriteString(fmt.Sprintf(`<input type="hidden" id="%s-id">`, target))
		// Find the model
		for _, m := range s.Program.Models {
			if lo(m.Name) == target {
				for _, f := range m.Fields {
					s.renderFormField(b, m, f)
				}
				break
			}
		}
		b.WriteString(`<button type="submit" class="btn btn-glow">` + svgIcon("check") + `<span>Salvar</span></button>`)
		b.WriteString(`</form></div>`)
	default:
		// Render children if any
		for _, child := range comp.Children {
			s.renderScreenComponent(b, child, primaryModel)
		}
	}
}

func (s *Servidor) renderChatComponent(b *strings.Builder, comp *ast.Component) {
	target := lo(comp.Target)
	messagesModel := lo(comp.Properties["messages_model"])
	if messagesModel == "" {
		messagesModel = "mensagem"
	}
	relationField := lo(comp.Properties["relation_field"])
	if relationField == "" {
		relationField = target
	}
	title := comp.Properties["title"]
	if title == "" {
		title = "Chat"
	}
	b.WriteString(fmt.Sprintf(`<div class="chat-shell card" data-chat-target="%s" data-chat-messages="%s" data-chat-relation="%s" data-chat-text="%s" data-chat-media="%s" data-chat-author="%s" data-chat-time="%s" data-chat-type="%s">`,
		target, messagesModel, relationField, lo(comp.Properties["text_field"]), lo(comp.Properties["media_field"]), lo(comp.Properties["author_field"]), lo(comp.Properties["timestamp_field"]), lo(comp.Properties["type_field"])))
	b.WriteString(`<div class="chat-side"><div class="chat-side-top"><h3>` + html.EscapeString(title) + `</h3><input class="chat-filter" placeholder="Buscar conversa" oninput="chatFilter('` + target + `',this.value)"></div>`)
	b.WriteString(`<div id="chat-conv-` + target + `" class="chat-conv-list"></div></div>`)
	b.WriteString(`<div class="chat-main"><div class="chat-main-top"><div><strong id="chat-title-` + target + `">Selecione uma conversa</strong><div id="chat-presence-` + target + `" class="chat-presence"></div></div><button class="btn btn-ghost btn-sm" type="button" onclick="refreshChat('` + target + `')">Atualizar</button></div>`)
	b.WriteString(`<div id="chat-msg-` + target + `" class="chat-messages"><div class="empty-state"><p>Nenhuma conversa selecionada</p></div></div>`)
	b.WriteString(`<div id="chat-typing-` + target + `" class="chat-typing"></div>`)
	b.WriteString(`<form class="chat-compose" onsubmit="chatSend('` + target + `',event)">`)
	b.WriteString(`<input type="file" id="chat-file-` + target + `" accept="image/*,audio/*,video/*,.pdf,.doc,.docx,.txt" onchange="chatUpload('` + target + `',this)">`)
	b.WriteString(`<input type="text" id="chat-input-` + target + `" placeholder="Digite uma mensagem" oninput="chatTyping('` + target + `',true)" onblur="chatTyping('` + target + `',false)">`)
	b.WriteString(`<button class="btn btn-glow" type="submit">Enviar</button>`)
	b.WriteString(`</form></div></div>`)
}

// ============================================================
// Model section (auto-generated CRUD)
// ============================================================

func (s *Servidor) renderModelSection(b *strings.Builder, model *ast.Model, theme *ast.Theme) {
	name := lo(model.Name)
	capName := cap(model.Name)

	b.WriteString(fmt.Sprintf(`<div id="secao-%s" class="section" style="display:none">`, name))

	// Section header
	b.WriteString(`<div class="sec-head">`)
	b.WriteString(`<div class="sec-left">`)
	b.WriteString(fmt.Sprintf(`<div class="sec-search"><input type="text" placeholder="Buscar em %s..." oninput="filtrar('%s',this.value)">`, html.EscapeString(capName), name))
	b.WriteString(svgIcon("search") + `</div></div>`)
	b.WriteString(`<div class="sec-actions">`)
	b.WriteString(fmt.Sprintf(`<button class="btn btn-ghost btn-sm" onclick="exportar('%s','csv')" title="CSV">%s<span>CSV</span></button>`, name, svgIcon("file")))
	b.WriteString(fmt.Sprintf(`<button class="btn btn-ghost btn-sm" onclick="exportar('%s','json')" title="JSON">%s<span>JSON</span></button>`, name, svgIcon("file")))
	b.WriteString(fmt.Sprintf(`<button class="btn btn-glow" onclick="abrirForm('%s')">%s<span>Novo %s</span></button>`, name, svgIcon("plus"), html.EscapeString(capName)))
	b.WriteString(`</div></div>`)

	// Table
	b.WriteString(`<div class="card table-wrap">`)
	b.WriteString(`<table><thead><tr><th class="th-id">#</th>`)
	for _, f := range model.Fields {
		if f.Type == ast.FieldSenha {
			continue
		}
		b.WriteString(`<th>` + html.EscapeString(cap(f.Name)) + `</th>`)
	}
	b.WriteString(`<th class="th-act"></th></tr></thead>`)
	b.WriteString(fmt.Sprintf(`<tbody id="tabela-%s"></tbody></table>`, name))
	b.WriteString(fmt.Sprintf(`<div id="paginacao-%s" class="pagination"></div>`, name))
	b.WriteString(fmt.Sprintf(`<div id="vazio-%s" class="empty-state">`, name))
	b.WriteString(svgIcon("inbox") + `<p>Nenhum registro</p></div></div>`)

	// Modal form
	b.WriteString(fmt.Sprintf(`<div id="modal-%s" class="modal-wrap" onclick="if(event.target===this)fecharForm('%s')">`, name, name))
	b.WriteString(`<div class="modal anim-modal">`)
	b.WriteString(fmt.Sprintf(`<div class="modal-top"><h2 id="titulo-form-%s">Novo %s</h2>`, name, capName))
	b.WriteString(fmt.Sprintf(`<button onclick="fecharForm('%s')" class="modal-x">`, name) + svgIcon("x") + `</button></div>`)
	b.WriteString(fmt.Sprintf(`<form onsubmit="salvar('%s',event)" class="modal-form"><input type="hidden" id="%s-id">`, name, name))

	for _, f := range model.Fields {
		s.renderFormField(b, model, f)
	}

	b.WriteString(`<div class="modal-foot">`)
	b.WriteString(`<button type="submit" class="btn btn-glow">` + svgIcon("check") + `<span>Salvar</span></button>`)
	b.WriteString(fmt.Sprintf(`<button type="button" class="btn btn-ghost" onclick="fecharForm('%s')">Cancelar</button>`, name))
	b.WriteString(`</div></form></div></div>`)
	b.WriteString(`</div>`) // section
}

// renderFormField generates the correct form input element based on field type.
func (s *Servidor) renderFormField(b *strings.Builder, model *ast.Model, f *ast.Field) {
	name := lo(model.Name)
	fname := lo(f.Name)
	req := ""
	if f.Required {
		req = " required"
	}

	b.WriteString(`<div class="field">`)
	b.WriteString(fmt.Sprintf(`<label for="%s-%s">%s</label>`, name, fname, html.EscapeString(cap(f.Name))))

	switch {
	// FK dropdown
	case f.Reference != "":
		refModel := lo(f.Reference)
		b.WriteString(fmt.Sprintf(`<select id="%s-%s" data-ref="%s"%s>`,
			name, fname, refModel, req))
		b.WriteString(`<option value="">Selecione...</option></select>`)

	// Enum dropdown
	case f.Type == ast.FieldEnum && len(f.EnumValues) > 0:
		b.WriteString(fmt.Sprintf(`<select id="%s-%s"%s>`, name, fname, req))
		b.WriteString(`<option value="">Selecione...</option>`)
		for _, v := range f.EnumValues {
			b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, html.EscapeString(v), html.EscapeString(cap(v))))
		}
		b.WriteString(`</select>`)

	// Status dropdown
	case f.Type == ast.FieldStatus:
		b.WriteString(fmt.Sprintf(`<select id="%s-%s"%s>`, name, fname, req))
		b.WriteString(`<option value="">Selecione...</option>`)
		for _, v := range []string{"ativo", "inativo", "pendente", "concluido"} {
			b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, v, cap(v)))
		}
		b.WriteString(`</select>`)

	// Long text
	case f.Type == ast.FieldTextoLongo:
		b.WriteString(fmt.Sprintf(`<textarea id="%s-%s" placeholder="%s" rows="4"%s></textarea>`,
			name, fname, html.EscapeString(cap(f.Name)), req))

	// File/image upload
	case f.Type == ast.FieldImagem || f.Type == ast.FieldUpload || f.Type == ast.FieldArquivo:
		b.WriteString(fmt.Sprintf(`<input type="hidden" id="%s-%s">`, name, fname))
		b.WriteString(fmt.Sprintf(`<input type="file" id="%s-%s-file" onchange="uploadFile('%s','%s',this)">`,
			name, fname, name, fname))
		b.WriteString(fmt.Sprintf(`<div id="%s-%s-preview" class="upload-preview"></div>`, name, fname))

	// Boolean checkbox
	case f.Type == ast.FieldBooleano:
		b.WriteString(fmt.Sprintf(`<label class="switch"><input type="checkbox" id="%s-%s"%s><span class="slider"></span></label>`,
			name, fname, req))

	// All other input types
	default:
		inputType := tipoInput(f.Type)
		extra := ""
		if f.Type == ast.FieldNumero || f.Type == ast.FieldDinheiro {
			extra = ` step="any"`
		}
		if f.Type == ast.FieldSenha {
			inputType = "password"
		}
		b.WriteString(fmt.Sprintf(`<input type="%s" id="%s-%s" placeholder="%s"%s%s>`,
			inputType, name, fname, html.EscapeString(cap(f.Name)), extra, req))
	}

	b.WriteString(`</div>`)
}

// ============================================================
// CSS Generation - fully theme-driven via CSS variables
// ============================================================

func (s *Servidor) generateCSS(theme *ast.Theme) string {
	// Compute derived colors for light/dark mode
	var darkBg, darkCard, darkText, darkText2, darkText3, darkBrd string
	var lightBg, lightCard, lightText, lightText2, lightText3, lightBrd string

	lightBg = "#f8fafc"
	lightCard = "rgba(255,255,255,0.85)"
	lightText = "#0f172a"
	lightText2 = "#64748b"
	lightText3 = "#94a3b8"
	lightBrd = "rgba(0,0,0,0.06)"

	darkBg = "#0c0a1d"
	darkCard = "rgba(30,27,75,0.6)"
	darkText = "#e2e8f0"
	darkText2 = "#94a3b8"
	darkText3 = "#64748b"
	darkBrd = "rgba(255,255,255,0.06)"

	// Use user overrides if dark mode is the default
	if theme.Dark {
		darkBg = theme.Background
		darkCard = theme.CardBg
		darkText = theme.TextColor
	} else {
		lightBg = theme.Background
		lightCard = theme.CardBg
		lightText = theme.TextColor
	}

	// Style-specific CSS
	styleCSS := s.styleVariantCSS(theme.Style)

	css := `
@import url('https://fonts.googleapis.com/css2?family=` + strings.ReplaceAll(theme.Font, " ", "+") + `:wght@300;400;500;600;700;800&display=swap');
*{margin:0;padding:0;box-sizing:border-box}

:root{
  --primary:` + theme.Primary + `;
  --secondary:` + theme.Secondary + `;
  --accent:` + theme.Accent + `;
  --sidebar-bg:` + theme.Sidebar + `;
  --radius:` + theme.Radius + `;
  --font:'` + theme.Font + `',system-ui,-apple-system,sans-serif;
  --bg:` + lightBg + `;
  --bg2:#f1f5f9;
  --card-bg:` + lightCard + `;
  --card-solid:#fff;
  --text:` + lightText + `;
  --text2:` + lightText2 + `;
  --text3:` + lightText3 + `;
  --border:` + lightBrd + `;
  --shadow:0 1px 2px rgba(0,0,0,.04),0 2px 8px rgba(0,0,0,.04);
  --shadow2:0 4px 24px rgba(0,0,0,.08);
  --shadow3:0 8px 40px rgba(0,0,0,.12);
  --ease:cubic-bezier(.4,0,.2,1);
  --dur:.25s;
}

body.dark{
  --bg:` + darkBg + `;
  --bg2:#12102a;
  --card-bg:` + darkCard + `;
  --card-solid:#1e1b4b;
  --text:` + darkText + `;
  --text2:` + darkText2 + `;
  --text3:` + darkText3 + `;
  --border:` + darkBrd + `;
  --shadow:0 1px 2px rgba(0,0,0,.2),0 2px 8px rgba(0,0,0,.2);
  --shadow2:0 4px 24px rgba(0,0,0,.4);
  --shadow3:0 8px 40px rgba(0,0,0,.5);
}

body{font-family:var(--font);background:radial-gradient(circle at top left,color-mix(in srgb,var(--primary) 14%,transparent),transparent 28%),
	radial-gradient(circle at top right,color-mix(in srgb,var(--accent) 12%,transparent),transparent 24%),
	linear-gradient(180deg,var(--bg),color-mix(in srgb,var(--bg) 88%,#031014));color:var(--text);
	display:flex;min-height:100vh;transition:background .4s var(--ease),color .3s var(--ease);overflow-x:hidden}
body::before{content:'';position:fixed;inset:0;background-image:linear-gradient(rgba(255,255,255,.018) 1px,transparent 1px),linear-gradient(90deg,rgba(255,255,255,.018) 1px,transparent 1px);background-size:28px 28px;pointer-events:none;opacity:.35}
body::after{content:'';position:fixed;inset:auto auto -140px -120px;width:360px;height:360px;border-radius:50%;background:color-mix(in srgb,var(--secondary) 18%,transparent);filter:blur(80px);pointer-events:none;opacity:.45}

/* ===== Sidebar ===== */
.sidebar{width:260px;background:linear-gradient(180deg,color-mix(in srgb,var(--sidebar-bg) 92%,#041116),color-mix(in srgb,var(--sidebar-bg) 78%,#02090b));color:#fff;display:flex;flex-direction:column;
	position:fixed;top:14px;left:14px;bottom:14px;z-index:50;transition:width .3s var(--ease),transform .3s var(--ease);border:1px solid rgba(255,255,255,.08);box-shadow:0 20px 50px rgba(0,0,0,.28), inset 0 1px 0 rgba(255,255,255,.04);border-radius:28px;overflow:hidden}
.sidebar.mini{width:72px}
.sb-top{flex:1;display:flex;flex-direction:column;overflow:hidden}
.sb-brand{padding:24px 18px 18px;display:flex;align-items:center;gap:12px;border-bottom:1px solid rgba(255,255,255,.08);background:linear-gradient(180deg,rgba(255,255,255,.05),transparent)}
.sb-logo{width:42px;height:42px;border-radius:14px;display:flex;align-items:center;justify-content:center;
	background:radial-gradient(circle at 30% 20%,color-mix(in srgb,var(--accent) 60%,#fff),var(--primary));flex-shrink:0;box-shadow:0 10px 24px color-mix(in srgb,var(--primary) 30%,transparent)}
.sb-logo svg{width:20px;height:20px}
.sb-logo img{border-radius:6px}
.sb-name{font-weight:800;font-size:1.1rem;white-space:nowrap;overflow:hidden;transition:opacity .2s;letter-spacing:-.03em}
.sb-collapse{margin-left:auto;background:none;border:none;color:rgba(255,255,255,.4);cursor:pointer;padding:4px;
  border-radius:6px;transition:all .2s;flex-shrink:0}
.sb-collapse:hover{color:#fff;background:rgba(255,255,255,.1)}
.sb-collapse svg{width:18px;height:18px;transition:transform .3s}
.sidebar.mini .sb-collapse svg{transform:rotate(180deg)}
.sidebar.mini .sb-name{opacity:0;width:0}
.sidebar.mini .sb-brand{justify-content:center;padding:20px 0}
.sidebar.mini .sb-collapse{display:none}

.sb-nav{padding:14px 10px;display:flex;flex-direction:column;gap:6px;flex:1;overflow-y:auto}
.sb-link{display:flex;align-items:center;gap:12px;padding:10px 12px;border-radius:calc(var(--radius) * 0.5);
	color:rgba(255,255,255,.6);text-decoration:none;font-size:.875rem;font-weight:600;
  transition:all .2s var(--ease);cursor:pointer;white-space:nowrap;position:relative;overflow:hidden}
.sb-link::before{content:'';position:absolute;inset:0;background:linear-gradient(90deg,rgba(255,255,255,.08),rgba(255,255,255,.02));opacity:0;transition:opacity .2s;border-radius:calc(var(--radius)*0.5)}
.sb-link:hover::before{opacity:1}
.sb-link:hover{color:rgba(255,255,255,.9)}
.sb-icon{width:36px;height:36px;display:flex;align-items:center;justify-content:center;border-radius:calc(var(--radius)*0.5);
	transition:background .2s,transform .2s;flex-shrink:0;background:rgba(255,255,255,.04)}
.sb-icon svg{width:18px;height:18px}
.sb-link.active{color:#fff}
.sb-link.active{background:linear-gradient(90deg,rgba(255,255,255,.08),rgba(255,255,255,.03));box-shadow:inset 0 1px 0 rgba(255,255,255,.05)}
.sb-link.active .sb-icon{background:linear-gradient(135deg,var(--primary),var(--secondary));box-shadow:0 8px 18px color-mix(in srgb,var(--primary) 40%,transparent);transform:translateY(-1px)}
.sidebar.mini .sb-link span{opacity:0;width:0}
.sidebar.mini .sb-nav{padding:12px 4px}
.sidebar.mini .sb-link{justify-content:center;padding:10px 0}

.sb-foot{padding:12px 16px;border-top:1px solid rgba(255,255,255,.08);display:flex;flex-direction:column;gap:8px}
.sb-theme{display:flex;align-items:center;gap:10px;background:none;border:none;color:rgba(255,255,255,.45);
  cursor:pointer;padding:8px;border-radius:calc(var(--radius)*0.5);font-size:.85rem;transition:all .2s;width:100%}
.sb-theme:hover{color:#fff;background:rgba(255,255,255,.08)}
.sb-theme svg{width:18px;height:18px;flex-shrink:0}
.sidebar.mini .sb-theme span{display:none}
.sidebar.mini .sb-theme{justify-content:center}
.sidebar.mini .sb-foot{align-items:center}
.sb-powered{font-size:.7rem;color:rgba(255,255,255,.2);text-align:center}
.sidebar.mini .sb-powered{display:none}

/* ===== Main ===== */
.main{margin-left:288px;flex:1;min-height:100vh;transition:margin-left .3s var(--ease)}
body.sb-mini .main{margin-left:72px}

/* ===== Topbar ===== */
.topbar{padding:14px 20px;display:flex;align-items:center;gap:16px;position:sticky;top:16px;z-index:30;
	margin:16px 20px 0;border:1px solid var(--border);background:color-mix(in srgb,var(--card-bg) 88%,transparent);border-radius:22px;transition:background .3s;
	backdrop-filter:blur(18px);-webkit-backdrop-filter:blur(18px);box-shadow:var(--shadow)}
.tb-title-wrap{display:flex;flex-direction:column;gap:2px;min-width:0}
.topbar h1{font-size:1.1rem;font-weight:800;letter-spacing:-.03em}
.tb-sub{font-size:.78rem;color:var(--text3);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.tb-menu{display:none;background:none;border:none;color:var(--text);cursor:pointer;padding:6px;border-radius:calc(var(--radius)*0.5)}
.tb-menu svg{width:22px;height:22px}
.tb-status{display:flex;align-items:center;gap:8px;flex:1;justify-content:center;min-width:0}
.tb-chip{display:inline-flex;align-items:center;gap:8px;padding:8px 12px;border-radius:999px;background:color-mix(in srgb,var(--bg2) 80%,transparent);border:1px solid var(--border);font-size:.78rem;font-weight:700;color:var(--text2);white-space:nowrap}
.tb-dot{width:8px;height:8px;border-radius:50%;display:inline-block;box-shadow:0 0 0 4px transparent}
.tb-dot-online{background:#22c55e;box-shadow:0 0 0 4px rgba(34,197,94,.14)}
.tb-dot-jobs{background:#f59e0b;box-shadow:0 0 0 4px rgba(245,158,11,.12)}
.tb-dot-wa{background:#06b6d4;box-shadow:0 0 0 4px rgba(6,182,212,.12)}
.tb-end{display:flex;align-items:center;gap:12px;margin-left:auto}
.tb-search{position:relative;display:flex;align-items:center;background:var(--bg);border:1px solid var(--border);border-radius:var(--radius);transition:border-color .2s,box-shadow .2s}
.tb-search:focus-within{border-color:var(--primary);box-shadow:0 0 0 3px color-mix(in srgb,var(--primary) 12%,transparent)}
.tb-search input{border:none;background:transparent;outline:none;font-size:.875rem;color:var(--text);width:200px;
  padding:8px 12px 8px 36px;transition:width .3s;font-family:var(--font)}
.tb-search input:focus{width:280px}
.tb-search svg{position:absolute;left:10px;width:16px;height:16px;color:var(--text3);pointer-events:none}

/* ===== Content ===== */
.content{padding:24px 28px 40px}

/* ===== Bento Grid ===== */
.bento{display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:18px;margin-bottom:28px}
.bento-card{position:relative;background:var(--card-bg);border:1px solid var(--border);border-radius:var(--radius);
	padding:24px;cursor:pointer;overflow:hidden;transition:all .3s var(--ease);box-shadow:var(--shadow);background-image:linear-gradient(180deg,rgba(255,255,255,.05),transparent)}
.bento-card:hover{transform:translateY(-6px) scale(1.01);box-shadow:var(--shadow2);border-color:color-mix(in srgb,var(--accent) 30%,var(--border))}
.bc-icon{width:52px;height:52px;border-radius:18px;display:flex;align-items:center;justify-content:center;
	background:linear-gradient(135deg,var(--accent),color-mix(in srgb,var(--accent) 70%,#fff));margin-bottom:16px;box-shadow:0 12px 28px color-mix(in srgb,var(--accent) 26%,transparent)}
.bc-icon svg{width:24px;height:24px;color:#fff}
.bc-num{font-size:clamp(1.75rem,3vw,2.25rem);font-weight:800;letter-spacing:-.03em;line-height:1}
.bc-label{font-size:.85rem;color:var(--text2);font-weight:600;margin-top:6px}
.bc-glow{position:absolute;top:-40%;right:-20%;width:120px;height:120px;border-radius:50%;
  background:var(--accent);opacity:.06;filter:blur(40px);pointer-events:none;transition:opacity .3s}
.bento-card:hover .bc-glow{opacity:.12}

/* ===== Dashboard Grid ===== */
.dash-grid{display:grid;grid-template-columns:2fr 1fr;gap:16px}
@media(max-width:900px){.dash-grid{grid-template-columns:1fr}}

/* ===== Card ===== */
.card{background:color-mix(in srgb,var(--card-bg) 92%,transparent);border:1px solid var(--border);border-radius:var(--radius);box-shadow:var(--shadow);
	overflow:hidden;transition:box-shadow .3s var(--ease),transform .3s var(--ease);backdrop-filter:blur(16px);-webkit-backdrop-filter:blur(16px)}
.card:hover{box-shadow:var(--shadow2);transform:translateY(-2px)}
.card-head{display:flex;align-items:center;gap:10px;padding:18px 20px;border-bottom:1px solid var(--border)}
.card-head svg{width:18px;height:18px;color:var(--primary)}
.card-head h3{font-size:.95rem;font-weight:600}

/* ===== Activity ===== */
.ativ-list{padding:8px 0;max-height:320px;overflow-y:auto}
.ativ-row{display:flex;align-items:center;gap:10px;padding:10px 20px;font-size:.875rem;transition:background .15s}
.ativ-row:hover{background:color-mix(in srgb,var(--primary) 4%,transparent)}
.ativ-tag{font-size:.7rem;padding:2px 8px;border-radius:99px;font-weight:700;color:#fff;text-transform:uppercase;letter-spacing:.5px;flex-shrink:0}
.ativ-c{background:#16a34a}.ativ-e{background:var(--primary)}.ativ-d{background:#dc2626}
.ativ-txt{flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.ativ-time{color:var(--text3);font-size:.8rem;font-variant-numeric:tabular-nums;flex-shrink:0}

/* ===== Info ===== */
.info-list{padding:4px 0}
.info-row{display:flex;justify-content:space-between;padding:12px 20px;border-bottom:1px solid var(--border);font-size:.875rem}
.info-row:last-child{border-bottom:none}
.info-k{color:var(--text2);font-weight:500}.info-v{font-weight:600}

/* ===== Charts ===== */
.chart-card{margin-bottom:20px}
.chart-wrap{padding:20px;min-height:120px}

/* ===== Section ===== */
.section{animation:fadeUp .35s var(--ease)}
@keyframes fadeUp{from{opacity:0;transform:translateY(12px)}to{opacity:1;transform:translateY(0)}}
.anim-in{animation:fadeUp .35s var(--ease)}
.sec-head{display:flex;align-items:center;justify-content:space-between;gap:14px;margin-bottom:22px;flex-wrap:wrap}
.sec-left{flex:1}
.sec-left h2{font-size:1.28rem;font-weight:800;letter-spacing:-.03em}
.sec-search{display:flex;align-items:center;max-width:380px;padding:0 14px;height:42px;background:var(--bg);
  border:1px solid var(--border);border-radius:var(--radius);transition:border-color .2s,box-shadow .2s}
.sec-search:focus-within{border-color:var(--primary);box-shadow:0 0 0 3px color-mix(in srgb,var(--primary) 12%,transparent)}
.sec-search input{flex:1;border:none;background:none;outline:none;font-size:.875rem;color:var(--text);padding:0 8px;font-family:var(--font)}
.sec-search input::placeholder{color:var(--text3)}
.sec-search svg{width:16px;height:16px;color:var(--text3);flex-shrink:0}
.sec-actions{display:flex;align-items:center;gap:8px}

/* ===== Table ===== */
.table-wrap{overflow-x:auto}
table{width:100%;border-collapse:collapse}
th{text-align:left;padding:12px 16px;font-weight:600;font-size:.75rem;text-transform:uppercase;
  letter-spacing:.6px;color:var(--text3);background:var(--bg2);border-bottom:1px solid var(--border)}
td{padding:13px 16px;border-bottom:1px solid var(--border);font-size:.875rem;transition:background .15s}
tr:hover td{background:color-mix(in srgb,var(--primary) 3%,transparent)}
.td-id{font-weight:700;color:var(--text3);font-size:.8rem;width:50px}
.th-id{width:50px}.th-act{width:90px;text-align:right}
.td-act{text-align:right;white-space:nowrap}
.row-anim{animation:fadeUp .25s var(--ease)}

/* Action btns */
.act-btn{width:34px;height:34px;display:inline-flex;align-items:center;justify-content:center;
  border:none;border-radius:calc(var(--radius)*0.5);cursor:pointer;transition:all .2s var(--ease);background:transparent}
.act-btn svg{width:15px;height:15px}
.act-edit{color:var(--primary)}.act-edit:hover{background:color-mix(in srgb,var(--primary) 10%,transparent);transform:scale(1.1)}
.act-del{color:#ef4444}.act-del:hover{background:rgba(239,68,68,.1);transform:scale(1.1)}

/* ===== Pagination ===== */
.pagination{display:flex;align-items:center;justify-content:center;gap:4px;padding:12px 16px;border-top:1px solid var(--border)}
.pagination button{min-width:34px;height:34px;display:inline-flex;align-items:center;justify-content:center;
  border:1px solid var(--border);border-radius:calc(var(--radius)*0.5);cursor:pointer;background:var(--bg);color:var(--text);
  font-size:.8rem;font-weight:600;transition:all .2s}
.pagination button:hover{border-color:var(--primary);color:var(--primary)}
.pagination button.active{background:var(--primary);color:#fff;border-color:var(--primary)}
.pagination button:disabled{opacity:.4;cursor:default}

/* ===== Empty state ===== */
.empty-state{display:flex;flex-direction:column;align-items:center;justify-content:center;padding:48px 20px;color:var(--text3);gap:8px}
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

.money{font-variant-numeric:tabular-nums;font-weight:600;color:var(--primary)}
.muted{color:var(--text3)}
.link{color:var(--primary);text-decoration:none;font-weight:500}
.link:hover{text-decoration:underline}

/* ===== Buttons ===== */
.btn{display:inline-flex;align-items:center;gap:7px;padding:10px 20px;border:none;border-radius:var(--radius);
  font-size:.875rem;font-weight:600;cursor:pointer;transition:all .25s var(--ease);text-decoration:none;
  position:relative;overflow:hidden;font-family:var(--font)}
.btn svg{width:16px;height:16px}
.btn-glow{background:linear-gradient(135deg,var(--primary),var(--secondary));color:#fff;
  box-shadow:0 2px 12px color-mix(in srgb,var(--primary) 35%,transparent)}
.btn-glow:hover{transform:translateY(-2px);box-shadow:0 6px 24px color-mix(in srgb,var(--primary) 45%,transparent)}
.btn-glow:active{transform:translateY(0)}
.btn-ghost{background:var(--bg2);color:var(--text2);border:1px solid var(--border)}
.btn-ghost:hover{background:var(--border);color:var(--text)}
.btn-sm{padding:7px 14px;font-size:.8rem}
.btn-sm svg{width:14px;height:14px}

/* ===== Modal ===== */
.modal-overlay{display:none;position:fixed;inset:0;background:rgba(0,0,0,.5);backdrop-filter:blur(8px);z-index:2000;align-items:center;justify-content:center}
.modal-overlay[style*="flex"]{display:flex!important}
.form-group{margin-bottom:14px}
.form-group label{display:block;font-size:.85rem;font-weight:500;margin-bottom:6px;opacity:.8}
.form-input{width:100%;padding:10px 14px;border-radius:calc(var(--radius)*0.6);border:1px solid var(--border);background:var(--bg2);color:var(--text);font-size:.9rem;outline:none;transition:border .2s}
.form-input:focus{border-color:var(--primary)}
.modal-wrap{display:none;position:fixed;inset:0;background:rgba(0,0,0,.4);backdrop-filter:blur(6px);
  -webkit-backdrop-filter:blur(6px);z-index:9000;align-items:center;justify-content:center;padding:20px}
.modal-wrap.show{display:flex}
.modal{width:100%;max-width:500px;max-height:85vh;overflow-y:auto;box-shadow:var(--shadow3);
  background:var(--card-bg);backdrop-filter:blur(24px);-webkit-backdrop-filter:blur(24px);
  border:1px solid var(--border);border-radius:var(--radius)}
.anim-modal{animation:modalIn .3s var(--ease)}
@keyframes modalIn{from{opacity:0;transform:scale(.95) translateY(10px)}to{opacity:1;transform:scale(1) translateY(0)}}
.modal-top{display:flex;align-items:center;justify-content:space-between;padding:20px 24px 0}
.modal-top h2{font-size:1.05rem;font-weight:700}
.modal-x{background:none;border:none;color:var(--text3);cursor:pointer;padding:6px;border-radius:calc(var(--radius)*0.5);transition:all .2s}
.modal-x:hover{background:var(--bg2);color:var(--text)}
.modal-x svg{width:18px;height:18px}
.modal-form{padding:16px 24px 24px}

/* ===== Form fields ===== */
.field{margin-bottom:16px}
.field label{display:block;font-weight:600;margin-bottom:6px;font-size:.8rem;color:var(--text2);text-transform:uppercase;letter-spacing:.5px}
.field input,.field select,.field textarea{width:100%;padding:11px 14px;border:1px solid var(--border);
  border-radius:calc(var(--radius)*0.5);font-size:.9rem;background:var(--bg);color:var(--text);
  transition:all .25s var(--ease);font-family:var(--font)}
.field input:focus,.field select:focus,.field textarea:focus{outline:none;border-color:var(--primary);
  box-shadow:0 0 0 4px color-mix(in srgb,var(--primary) 10%,transparent);background:var(--card-solid)}
.field input::placeholder,.field textarea::placeholder{color:var(--text3)}
.field textarea{resize:vertical;min-height:80px}
.field input[type="file"]{padding:8px;cursor:pointer}
.upload-preview{min-height:0}
.upload-preview img{display:block}
.modal-foot{display:flex;gap:10px;padding-top:16px;border-top:1px solid var(--border)}

/* ===== Toggle switch ===== */
.switch{position:relative;display:inline-block;width:44px;height:24px;cursor:pointer}
.switch input{opacity:0;width:0;height:0}
.slider{position:absolute;inset:0;background:var(--border);border-radius:24px;transition:.3s}
.slider::before{content:'';position:absolute;width:18px;height:18px;left:3px;bottom:3px;background:#fff;border-radius:50%;transition:.3s}
.switch input:checked+.slider{background:var(--primary)}
.switch input:checked+.slider::before{transform:translateX(20px)}

/* ===== Toast ===== */
.toast{position:fixed;bottom:24px;right:24px;padding:14px 28px;border-radius:var(--radius);color:#fff;
  font-weight:600;font-size:.9rem;z-index:200;opacity:0;transform:translateY(16px) scale(.95);
  transition:all .35s var(--ease);pointer-events:none;backdrop-filter:blur(8px);font-family:var(--font)}
.toast.show{opacity:1;transform:translateY(0) scale(1)}
.toast.ok{background:linear-gradient(135deg,#16a34a,#15803d);box-shadow:0 4px 20px rgba(22,163,74,.35)}
.toast.erro{background:linear-gradient(135deg,#ef4444,#dc2626);box-shadow:0 4px 20px rgba(239,68,68,.35)}

/* ===== Scrollbar ===== */
::-webkit-scrollbar{width:6px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--border);border-radius:3px}
::-webkit-scrollbar-thumb:hover{background:var(--text3)}

/* ===== Responsive ===== */
.chat-shell{display:grid;grid-template-columns:340px 1fr;min-height:720px;overflow:hidden;background:linear-gradient(180deg,color-mix(in srgb,var(--card-bg) 95%,transparent),color-mix(in srgb,var(--bg2) 18%,transparent))}
.chat-side{border-right:1px solid var(--border);display:flex;flex-direction:column;background:linear-gradient(180deg,color-mix(in srgb,var(--bg2) 66%,transparent),transparent)}
.chat-side-top{padding:18px;border-bottom:1px solid var(--border);display:flex;flex-direction:column;gap:10px;background:linear-gradient(180deg,rgba(255,255,255,.04),transparent)}
.chat-filter{width:100%;padding:10px 12px;border:1px solid var(--border);border-radius:calc(var(--radius)*0.5);background:var(--bg);color:var(--text)}
.chat-conv-list{overflow:auto;display:flex;flex-direction:column}
.chat-conv{padding:14px 16px;border-bottom:1px solid color-mix(in srgb,var(--border) 65%,transparent);cursor:pointer;transition:background .2s,transform .2s;display:grid;grid-template-columns:44px 1fr auto;gap:12px;align-items:center}
.chat-conv:hover,.chat-conv.active{background:linear-gradient(90deg,color-mix(in srgb,var(--primary) 10%,transparent),transparent);transform:translateX(2px)}
.chat-conv-avatar{width:44px;height:44px;border-radius:14px;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,color-mix(in srgb,var(--secondary) 70%,#fff),var(--primary));color:#fff;font-weight:800;letter-spacing:-.03em;box-shadow:0 10px 24px color-mix(in srgb,var(--primary) 18%,transparent)}
.chat-conv-body{min-width:0}
.chat-conv-name{font-weight:700;font-size:.92rem}
.chat-conv-last{font-size:.82rem;color:var(--text2);margin-top:4px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.chat-conv-meta{display:flex;justify-content:space-between;gap:8px;margin-top:6px;font-size:.75rem;color:var(--text3)}
.chat-unread{min-width:24px;height:24px;padding:0 8px;border-radius:999px;background:linear-gradient(135deg,var(--accent),color-mix(in srgb,var(--accent) 70%,#fff));color:#042225;display:inline-flex;align-items:center;justify-content:center;font-size:.72rem;font-weight:800;box-shadow:0 8px 18px color-mix(in srgb,var(--accent) 18%,transparent)}
.chat-main{display:flex;flex-direction:column;min-width:0;background:linear-gradient(180deg,transparent,color-mix(in srgb,var(--bg2) 14%,transparent))}
.chat-main-top{padding:18px 20px;border-bottom:1px solid var(--border);display:flex;justify-content:space-between;align-items:center;gap:12px;background:linear-gradient(180deg,rgba(255,255,255,.04),transparent)}
.chat-presence{font-size:.8rem;color:var(--text2);margin-top:4px}
.chat-messages{flex:1;padding:24px;background:radial-gradient(circle at top left,color-mix(in srgb,var(--primary) 8%,transparent),transparent 28%),linear-gradient(180deg,color-mix(in srgb,var(--bg2) 78%,transparent),transparent);overflow:auto;display:flex;flex-direction:column;gap:12px}
.chat-bubble{max-width:min(72%,680px);padding:13px 15px;border-radius:18px;background:var(--bg);border:1px solid var(--border);box-shadow:var(--shadow);position:relative}
.chat-bubble.mine{align-self:flex-end;background:linear-gradient(135deg,color-mix(in srgb,var(--primary) 18%,var(--card-solid)),color-mix(in srgb,var(--secondary) 14%,var(--card-solid)));border-bottom-right-radius:8px}
.chat-bubble.other{align-self:flex-start}
.chat-bubble-meta{font-size:.72rem;color:var(--text3);margin-top:6px}
.chat-media{margin-top:8px}
.chat-media img,.chat-media video{max-width:100%;border-radius:12px}
.chat-media audio{width:100%}
.chat-compose{display:flex;gap:10px;padding:14px 16px;border-top:1px solid var(--border);align-items:center;background:linear-gradient(180deg,transparent,color-mix(in srgb,var(--bg2) 18%,transparent))}
.chat-compose input[type="file"]{max-width:180px}
.chat-compose input[type="text"]{flex:1;padding:13px 16px;border:1px solid var(--border);border-radius:999px;background:color-mix(in srgb,var(--bg) 88%,transparent);color:var(--text)}
.chat-typing{min-height:20px;padding:0 18px 8px;color:var(--text2);font-size:.78rem}
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
	.tb-status{display:none}
  .modal{max-width:95vw}
  .dash-grid{grid-template-columns:1fr}
	.chat-shell{grid-template-columns:1fr;min-height:unset}
	.chat-side{max-height:240px}
}
@media(max-width:480px){
  .bento{grid-template-columns:1fr}
  .content{padding:12px}
}

` + styleCSS

	// Inject user custom CSS
	if theme.CustomCSS != "" {
		css += "\n/* === User Custom CSS === */\n" + theme.CustomCSS + "\n"
	}

	return css
}

// styleVariantCSS returns extra CSS for the chosen style variant.
func (s *Servidor) styleVariantCSS(style string) string {
	switch style {
	case "flat":
		return `
/* === Flat Style === */
.card,.modal,.topbar,.bento-card{backdrop-filter:none;-webkit-backdrop-filter:none;box-shadow:none!important}
.card{background:var(--card-solid);border:1px solid var(--border)}
.modal{background:var(--card-solid)}
.topbar{background:var(--card-solid);backdrop-filter:none}
.bento-card{background:var(--card-solid)}
.bento-card:hover{transform:none;box-shadow:none!important}
.bc-glow{display:none}
.btn-glow{box-shadow:none}
.btn-glow:hover{transform:none;box-shadow:none}
`
	case "neumorphism":
		return `
/* === Neumorphism Style === */
.card,.bento-card{backdrop-filter:none;-webkit-backdrop-filter:none;background:var(--bg);border:none;
  box-shadow:6px 6px 12px color-mix(in srgb,var(--bg) 80%,#000),
             -6px -6px 12px color-mix(in srgb,var(--bg) 80%,#fff)!important}
.card:hover,.bento-card:hover{box-shadow:8px 8px 16px color-mix(in srgb,var(--bg) 75%,#000),
             -8px -8px 16px color-mix(in srgb,var(--bg) 75%,#fff)!important}
body.dark .card,body.dark .bento-card{
  box-shadow:6px 6px 12px rgba(0,0,0,.4),-6px -6px 12px rgba(255,255,255,.03)!important}
body.dark .card:hover,body.dark .bento-card:hover{
  box-shadow:8px 8px 16px rgba(0,0,0,.5),-8px -8px 16px rgba(255,255,255,.04)!important}
.modal{backdrop-filter:none;background:var(--bg);border:none;
  box-shadow:8px 8px 20px color-mix(in srgb,var(--bg) 75%,#000),
             -8px -8px 20px color-mix(in srgb,var(--bg) 75%,#fff)!important}
.bc-glow{display:none}
.bento-card:hover{transform:translateY(-2px)}
.field input,.field select,.field textarea{border:none;
  box-shadow:inset 3px 3px 6px color-mix(in srgb,var(--bg) 85%,#000),
             inset -3px -3px 6px color-mix(in srgb,var(--bg) 85%,#fff)}
.field input:focus,.field select:focus,.field textarea:focus{
  box-shadow:inset 3px 3px 6px color-mix(in srgb,var(--bg) 85%,#000),
             inset -3px -3px 6px color-mix(in srgb,var(--bg) 85%,#fff),
             0 0 0 3px color-mix(in srgb,var(--primary) 15%,transparent)}
`
	case "minimal":
		return `
/* === Minimal Style === */
.card,.bento-card,.modal{backdrop-filter:none;-webkit-backdrop-filter:none;background:transparent;
  box-shadow:none!important;border:1px solid var(--border);border-radius:calc(var(--radius)*0.5)}
.topbar{background:transparent;backdrop-filter:none;border-bottom:1px solid var(--border)}
.bento-card{padding:20px}
.bento-card:hover{transform:none;border-color:var(--primary)}
.bc-glow{display:none}
.bc-icon{border-radius:8px;width:40px;height:40px}
.btn-glow{box-shadow:none;border-radius:calc(var(--radius)*0.3)}
.btn-glow:hover{transform:none;box-shadow:0 2px 8px color-mix(in srgb,var(--primary) 20%,transparent)}
.card-head{padding:14px 16px}
.content{padding:32px 40px}
@media(max-width:768px){.content{padding:16px}}
`
	default: // glassmorphism (default)
		return `
/* === Glassmorphism Style === */
.card{backdrop-filter:blur(16px);-webkit-backdrop-filter:blur(16px)}
.modal{backdrop-filter:blur(24px);-webkit-backdrop-filter:blur(24px)}
.topbar{backdrop-filter:blur(16px);-webkit-backdrop-filter:blur(16px)}
`
	}
}

// ============================================================
// JavaScript Generation
// ============================================================

func (s *Servidor) generateJS(theme *ast.Theme) string {
	var b strings.Builder

	// Icon strings for table action buttons
	editIcon := strings.ReplaceAll(strings.ReplaceAll(svgIcon("edit"), `"`, `'`), "\n", "")
	trashIcon := strings.ReplaceAll(strings.ReplaceAll(svgIcon("trash"), `"`, `'`), "\n", "")
	b.WriteString(fmt.Sprintf("var ICO_E=%q,ICO_D=%q;\n", editIcon, trashIcon))

	// Model metadata for JS
	b.WriteString("var M={\n")
	for _, model := range s.Program.Models {
		name := lo(model.Name)
		b.WriteString(fmt.Sprintf("'%s':[", name))
		for i, f := range model.Fields {
			ft := fieldTypeCode(f)
			ref := ""
			if f.Reference != "" {
				ref = fmt.Sprintf(",r:'%s'", lo(f.Reference))
			}
			enumVals := ""
			if f.Type == ast.FieldEnum && len(f.EnumValues) > 0 {
				enumVals = fmt.Sprintf(",ev:%s", enumJSArray(f.EnumValues))
			}
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(fmt.Sprintf("{n:'%s',t:'%s'%s%s}", lo(f.Name), ft, ref, enumVals))
		}
		b.WriteString("],\n")
	}
	b.WriteString("};\n")

	// Pagination state
	b.WriteString("var PAGE_SIZE=20,pages={};\n")

	// Chart.js theme colors
	b.WriteString(fmt.Sprintf("var THEME_PRI='%s',THEME_SEC='%s',THEME_ACC='%s';\n", theme.Primary, theme.Secondary, theme.Accent))
	b.WriteString("var chartColors=[THEME_PRI,THEME_SEC,THEME_ACC,'#10b981','#3b82f6','#ef4444','#06b6d4','#ec4899'];\n")

	b.WriteString(`
var ativs=[];
function $(id){return document.getElementById(id);}
function esc(v){if(v==null)return'';var d=document.createElement('div');d.textContent=String(v);return d.innerHTML;}

// ===== Navigation =====
function irPara(n,el){
  document.querySelectorAll('.section').forEach(function(s){s.style.display='none';});
  var sec=$('secao-'+n);
  if(sec){sec.style.display='block';sec.classList.add('anim-in');}
  document.querySelectorAll('.sb-link').forEach(function(a){a.classList.remove('active');});
  if(el)el.classList.add('active');
  var title=n==='dashboard'?'Dashboard':n.replace('screen-','').charAt(0).toUpperCase()+n.replace('screen-','').slice(1);
  $('page-title').textContent=title;
  if(innerWidth<768)$('sidebar').classList.remove('open');
  // Load data for inline lists in custom screens
  if(sec){carregarListasInline();}
  // Also load data for model section
  if(n!=='dashboard'&&!n.startsWith('screen-')){carregar(n);}
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

// ===== Form open/close =====
function abrirForm(m){
  $('modal-'+m).classList.add('show');
  $(m+'-id').value='';
  $('modal-'+m).querySelector('form').reset();
  $('titulo-form-'+m).textContent='Novo '+m.charAt(0).toUpperCase()+m.slice(1);
  M[m].forEach(function(c){
    if(c.t==='f'){var prev=$(m+'-'+c.n+'-preview');if(prev)prev.innerHTML='';}
  });
  carregarSelects(m);
}
function fecharForm(m){$('modal-'+m).classList.remove('show');}

// ===== Search / Filter =====
function filtrar(m,q){
  q=q.toLowerCase();
  document.querySelectorAll('#tabela-'+m+' tr').forEach(function(r){
    r.style.display=r.textContent.toLowerCase().includes(q)?'':'none';
  });
}
function buscaGlobal(q){
  if(!q){document.querySelectorAll('table tr').forEach(function(r){r.style.display='';});return;}
  q=q.toLowerCase();
  document.querySelectorAll('table tbody tr').forEach(function(r){r.style.display=r.textContent.toLowerCase().includes(q)?'':'none';});
}

// ===== File upload =====
function uploadFile(m,fname,input){
  if(!input.files||!input.files[0])return;
  var fd=new FormData();fd.append('file',input.files[0]);
  var prev=$(m+'-'+fname+'-preview');
  if(prev)prev.innerHTML='<span style="color:var(--text2);font-size:.85rem">Enviando...</span>';
  fetch('/upload',{method:'POST',body:fd})
    .then(function(r){if(!r.ok)throw new Error('Upload falhou');return r.json();})
    .then(function(d){
      $(m+'-'+fname).value=d.path;
      if(prev){
        if(d.path.match(/\.(jpg|jpeg|png|gif|webp|svg)$/i)){
          prev.innerHTML='<img src="'+esc(d.path)+'" style="max-width:100%;max-height:120px;border-radius:8px;margin-top:6px">';
        }else{
          prev.innerHTML='<span style="color:var(--primary);font-size:.85rem;margin-top:4px;display:block">'+esc(d.name)+' &#10003;</span>';
        }
      }
    })
    .catch(function(err){toast('Erro upload: '+err.message,'erro');if(prev)prev.innerHTML='';});
}

// ===== FK & Enum select population =====
function carregarSelects(m){
  M[m].forEach(function(c){
    if(!c.r)return;
    var sel=$(m+'-'+c.n);if(!sel||sel.tagName!=='SELECT')return;
    fetch('/api/'+c.r).then(function(r){return r.json();}).then(function(items){
      var val=sel.value;
      sel.innerHTML='<option value="">Selecione...</option>';
      if(!items||!items.length)return;
      var labelKey=null;
      if(M[c.r]){for(var i=0;i<M[c.r].length;i++){if(M[c.r][i].t==='t'||M[c.r][i].t==='e'){labelKey=M[c.r][i].n;break;}}}
      items.forEach(function(it){
        var label=labelKey?it[labelKey]:(it.nome||it.name||it.titulo||it.title||'#'+it.id);
        var o=document.createElement('option');o.value=it.id;o.textContent=label;
        if(String(it.id)===String(val))o.selected=true;
        sel.appendChild(o);
      });
    });
  });
}

// ===== Table cell formatting =====
function fmtCell(v,t){
  var s=esc(v);
  if(!s||s==='-')return'<span class="muted">&#8212;</span>';
  if(t==='s')return'<span class="pill pill-'+pillColor(v)+'">'+s+'</span>';
  if(t==='d'){var n=parseFloat(v);return'<span class="money">R$&nbsp;'+n.toFixed(2)+'</span>';}
  if(t==='e')return'<a class="link" href="mailto:'+s+'">'+s+'</a>';
  if(t==='en')return'<span class="pill pill-blue">'+s+'</span>';
  if(t==='b'){return v?'<span class="pill pill-green">Sim</span>':'<span class="pill pill-red">N&atilde;o</span>';}
  if(t==='f'){
    if(String(v).match(/\.(jpg|jpeg|png|gif|webp|svg)$/i))return'<img src="'+s+'" style="max-height:40px;border-radius:4px">';
    return'<a class="link" href="'+s+'" target="_blank">'+s.split('/').pop()+'</a>';
  }
  if(t==='tl'){return s.length>60?s.substring(0,60)+'...':s;}
  return s;
}

function pillColor(v){
  if(!v)return'blue';v=v.toLowerCase();
  if('ativo,livre,aberto,ok,sim,disponivel,pronto,entregue,pago,aprovado,online,concluido'.indexOf(v)>=0)return'green';
  if('inativo,ocupado,fechado,nao,cancelado,bloqueado,offline,reprovado'.indexOf(v)>=0)return'red';
  if('pendente,aguardando,em andamento,reservado,preparando,analise'.indexOf(v)>=0)return'yellow';
  return'blue';
}

// ===== Activity feed =====
function addAtiv(tipo,mod,nome){
  var labs={c:'Criado',e:'Editado',d:'Exclu\u00eddo'};
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
    if(a.n)h+=' \u2014 '+esc(a.n);
    h+='</span><span class="ativ-time">'+a.h+'</span></div>';
  });
  el.innerHTML=h;
}

function renderTableRows(tb,m,items){
	var cs=M[m];
	if(!tb||!cs)return;
	tb.innerHTML='';
	(items||[]).forEach(function(item){
		var tr=document.createElement('tr');tr.className='row-anim';
		var h='<td class="td-id">'+item.id+'</td>';
		cs.forEach(function(c){
			if(c.t==='pw')return;
			h+='<td>'+fmtCell(item[c.n],c.t)+'</td>';
		});
		h+='<td class="td-act"><button class="act-btn act-edit" onclick="editar(\''+m+'\','+item.id+')">'+ICO_E+'</button>';
		h+='<button class="act-btn act-del" onclick="excluir(\''+m+'\','+item.id+')">'+ICO_D+'</button></td>';
		tr.innerHTML=h;tb.appendChild(tr);
	});
}

function carregarListasInline(modelo){
	var selector=modelo?'[data-list-model="'+modelo+'"]':'[data-list-model]';
	document.querySelectorAll(selector).forEach(function(card){
		var m=card.getAttribute('data-list-model');
		var tb=card.querySelector('tbody');
		var vazio=card.querySelector('.empty-state');
		var table=card.querySelector('table');
		if(!m||!tb||!table)return;
		fetch('/api/'+m).then(function(r){return r.json();}).then(function(items){
			items=items||[];
			renderTableRows(tb,m,items);
			if(!items.length){
				if(vazio)vazio.style.display='flex';
				table.style.display='none';
				return;
			}
			if(vazio)vazio.style.display='none';
			table.style.display='';
		}).catch(function(){
			if(vazio)vazio.style.display='flex';
			table.style.display='none';
		});
	});
}

// ===== Data loading with pagination =====
function carregar(m,page){
  if(!page)page=1;
  pages[m]=page;
  fetch('/api/'+m).then(function(r){return r.json();}).then(function(items){
    var tb=$('tabela-'+m),vz=$('vazio-'+m),st=$('stat-'+m);
    if(!tb)return;
    tb.innerHTML='';
    var total=items?items.length:0;
    if(st)st.textContent=total;
    if(!items||!items.length){
      vz.style.display='flex';tb.closest('table').style.display='none';
      var pg=$('paginacao-'+m);if(pg)pg.innerHTML='';
      return;
    }
    vz.style.display='none';tb.closest('table').style.display='';

    // Pagination
    var totalPages=Math.ceil(total/PAGE_SIZE);
    var start=(page-1)*PAGE_SIZE;
    var end=Math.min(start+PAGE_SIZE,total);
    var pageItems=items.slice(start,end);

    var cs=M[m];
		renderTableRows(tb,m,pageItems);

    // Render pagination controls
    var pg=$('paginacao-'+m);
    if(pg&&totalPages>1){
      var ph='<button '+(page<=1?'disabled':'')+' onclick="carregar(\''+m+'\','+(page-1)+')">&laquo;</button>';
      for(var i=1;i<=totalPages;i++){
        if(totalPages>7&&Math.abs(i-page)>2&&i!==1&&i!==totalPages){
          if(i===2||i===totalPages-1)ph+='<button disabled>...</button>';
          continue;
        }
        ph+='<button class="'+(i===page?'active':'')+'" onclick="carregar(\''+m+'\','+i+')">'+i+'</button>';
      }
      ph+='<button '+(page>=totalPages?'disabled':'')+' onclick="carregar(\''+m+'\','+(page+1)+')">&raquo;</button>';
      pg.innerHTML=ph;
    }else if(pg){pg.innerHTML='';}
  });
}

// ===== CRUD operations =====
function salvar(m,e){
  e.preventDefault();var id=$(m+'-id').value;var d={};
  M[m].forEach(function(c){
    var el=$(m+'-'+c.n);if(!el)return;
    if(c.t==='b'){d[c.n]=el.checked;return;}
    var v=el.value;
    d[c.n]=(c.t==='n'||c.t==='d')?parseFloat(v)||0:v;
  });
  fetch(id?'/api/'+m+'/'+id:'/api/'+m,{method:id?'PUT':'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(d)})
    .then(function(r){if(!r.ok)return r.json().then(function(e){throw new Error(e.erro||e.error||'Erro');});return r.json();})
		.then(function(){fecharForm(m);carregar(m);carregarListasInline(m);addAtiv(id?'e':'c',m,d[M[m][0].n]||'');toast(id?'Atualizado!':'Criado!');renderCharts();})
    .catch(function(err){toast('Erro: '+err.message,'erro');});
}

function editar(m,id){
  fetch('/api/'+m+'/'+id).then(function(r){return r.json();}).then(function(item){
    $(m+'-id').value=item.id;
    carregarSelects(m);
    M[m].forEach(function(c){
      var el=$(m+'-'+c.n);
      if(!el)return;
      if(c.t==='b'){el.checked=!!item[c.n];return;}
      el.value=item[c.n]||'';
      if(c.t==='f'){
        var prev=$(m+'-'+c.n+'-preview');
        if(prev&&item[c.n]){
          if(String(item[c.n]).match(/\.(jpg|jpeg|png|gif|webp|svg)$/i)){
            prev.innerHTML='<img src="'+esc(item[c.n])+'" style="max-width:100%;max-height:120px;border-radius:8px;margin-top:6px">';
          }else{
            prev.innerHTML='<span style="color:var(--primary);font-size:.85rem;margin-top:4px;display:block">'+esc(item[c.n])+'</span>';
          }
        }
      }
      if(c.r&&el){setTimeout(function(){el.value=item[c.n]||'';},300);}
    });
    $('titulo-form-'+m).textContent='Editar';
    $('modal-'+m).classList.add('show');
  });
}

function excluir(m,id){
  if(!confirm('Excluir #'+id+'?'))return;
  var tb=$('tabela-'+m),rows=tb.querySelectorAll('tr'),label='';
  rows.forEach(function(r){if(r.querySelector('.td-id')&&r.querySelector('.td-id').textContent==id){label=r.children[1]?r.children[1].textContent:'';}});
	fetch('/api/'+m+'/'+id,{method:'DELETE'}).then(function(){carregar(m);carregarListasInline(m);addAtiv('d',m,label);toast('Excluido!');renderCharts();});
}

function exportar(m,fmt){
  var a=document.createElement('a');a.href='/api/'+m+'/export/'+fmt;a.download='';document.body.appendChild(a);a.click();document.body.removeChild(a);
}

// ===== Chart.js rendering =====
var chartInstances={};
var CHATS={};

function refreshSystemStatus(){
	fetch('/api/_jobs/status').then(function(r){return r.json();}).then(function(d){
		var queued=(d&&d.queued)||0, running=(d&&d.running)||0;
		if($('tb-jobs')) $('tb-jobs').textContent='jobs '+queued+'/'+running;
	}).catch(function(){});
	fetch('/api/whatsapp/sessions').then(function(r){if(!r.ok) throw new Error(); return r.json();}).then(function(items){
		var connected=(items||[]).filter(function(it){return !!it.connected;}).length;
		if($('tb-wa')) $('tb-wa').textContent=connected>0?('whatsapp '+connected+' online'):('whatsapp offline');
	}).catch(function(){ if($('tb-wa')) $('tb-wa').textContent='whatsapp offline'; });
	if($('tb-sockets') && typeof ws!=='undefined' && ws){ $('tb-sockets').textContent='tempo real ativo'; }
}

function renderCharts(){
  fetch('/api/_stats').then(function(r){return r.json();}).then(function(stats){
    var models=Object.keys(stats);

    // Bar chart - records per model
    var el=$('chart-models');
    if(el){
      var labels=models.map(function(m){return m.charAt(0).toUpperCase()+m.slice(1);});
      var data=models.map(function(m){return stats[m].count||0;});
      var bgColors=models.map(function(_,i){return chartColors[i%chartColors.length];});
      if(chartInstances['models'])chartInstances['models'].destroy();
      chartInstances['models']=new Chart(el,{
        type:'bar',
        data:{labels:labels,datasets:[{label:'Registros',data:data,backgroundColor:bgColors,borderRadius:6,borderSkipped:false}]},
        options:{responsive:true,maintainAspectRatio:false,plugins:{legend:{display:false}},
          scales:{y:{beginAtZero:true,grid:{color:'rgba(128,128,128,0.1)'},ticks:{color:getComputedStyle(document.body).getPropertyValue('--text2').trim()}},
                  x:{grid:{display:false},ticks:{color:getComputedStyle(document.body).getPropertyValue('--text2').trim()}}}}
      });
    }

    // Status doughnut chart
    var sel=$('chart-status');
    if(sel){
      var statusData={};
      models.forEach(function(m){
        var st=stats[m].statuses;
        if(!st)return;
        Object.keys(st).forEach(function(k){statusData[k]=(statusData[k]||0)+st[k];});
      });
      var sKeys=Object.keys(statusData);
      if(sKeys.length){
        var sColors=sKeys.map(function(k){var pc=pillColor(k);return{green:'#16a34a',red:'#ef4444',yellow:'#f59e0b',blue:'#3b82f6'}[pc]||THEME_PRI;});
        if(chartInstances['status'])chartInstances['status'].destroy();
        chartInstances['status']=new Chart(sel,{
          type:'doughnut',
          data:{labels:sKeys.map(function(k){return k.charAt(0).toUpperCase()+k.slice(1);}),
                datasets:[{data:sKeys.map(function(k){return statusData[k];}),backgroundColor:sColors,borderWidth:0}]},
          options:{responsive:true,maintainAspectRatio:false,cutout:'60%',
            plugins:{legend:{position:'bottom',labels:{color:getComputedStyle(document.body).getPropertyValue('--text').trim(),padding:16,usePointStyle:true,pointStyle:'circle'}}}}
        });
      }
    }

    // Custom chart canvases
    document.querySelectorAll('canvas[data-chart-model]').forEach(function(canvas){
      var cid=canvas.id;
      var ctype=canvas.getAttribute('data-chart-type')||'bar';
      var cmodel=canvas.getAttribute('data-chart-model');
      if(!stats[cmodel])return;
      fetch('/api/'+cmodel).then(function(r){return r.json();}).then(function(items){
        if(!items||!items.length)return;
        if(chartInstances[cid])chartInstances[cid].destroy();
        // Try to find a numeric field and a label field
        var cs=M[cmodel];if(!cs)return;
        var numField=null,labelField=null;
        cs.forEach(function(c){if(!numField&&(c.t==='n'||c.t==='d'))numField=c.n;if(!labelField&&c.t==='t')labelField=c.n;});
        if(!numField)return;
        if(!labelField)labelField=cs[0].n;
        var labels=items.map(function(it){return it[labelField]||'#'+it.id;});
        var data=items.map(function(it){return parseFloat(it[numField])||0;});
        chartInstances[cid]=new Chart(canvas,{
          type:ctype,
          data:{labels:labels,datasets:[{label:numField,data:data,
            backgroundColor:items.map(function(_,i){return chartColors[i%chartColors.length];}),
            borderColor:ctype==='line'?THEME_PRI:undefined,
            borderWidth:ctype==='line'?2:0,borderRadius:ctype==='bar'?6:0,
            fill:ctype==='line'?false:undefined,tension:0.3}]},
          options:{responsive:true,maintainAspectRatio:false,
            plugins:{legend:{display:false}},
            scales:ctype==='doughnut'||ctype==='pie'?{}:{
              y:{beginAtZero:true,grid:{color:'rgba(128,128,128,0.1)'}},
              x:{grid:{display:false}}}}
        });
      });
    });
  }).catch(function(){});
}

// ===== WebSocket for real-time updates =====
var ws;
function connectWS(){
  var proto=location.protocol==='https:'?'wss:':'ws:';
  ws=new WebSocket(proto+'//'+location.host+'/ws');
  ws.onmessage=function(e){
    try{
      var msg=JSON.parse(e.data);
			if(msg.model){carregar(msg.model);carregarListasInline(msg.model);renderCharts();chatHandleUpdate(msg);}
			if(msg.type==='presenca'||msg.type==='digitando'||msg.type==='presenca_socket'){chatHandlePresence(msg);}
			if(msg.type==='qr'){toast('QR code atualizado para sessao '+(msg.session||'default'));}
			if(msg.type==='whatsapp_status'&&msg.data&&msg.data.status){toast('WhatsApp '+msg.data.status);}
			if(msg.type==='presenca_socket'&&$('tb-sockets')&&msg.data){$('tb-sockets').textContent=(msg.data.connections||0)+' conexoes';}
			if(msg.type==='whatsapp_status'){refreshSystemStatus();}
    }catch(ex){}
  };
  ws.onclose=function(){setTimeout(connectWS,2000);};
  ws.onerror=function(){ws.close();};
}

function initChats(){
	document.querySelectorAll('[data-chat-target]').forEach(function(el){
		var target=el.getAttribute('data-chat-target');
		CHATS[target]={
			target:target,
			messages:el.getAttribute('data-chat-messages')||'mensagem',
			relation:el.getAttribute('data-chat-relation')||target,
			textField:el.getAttribute('data-chat-text')||'corpo',
			mediaField:el.getAttribute('data-chat-media')||'media_url',
			authorField:el.getAttribute('data-chat-author')||'de_mim',
			timeField:el.getAttribute('data-chat-time')||'criado_em',
			typeField:el.getAttribute('data-chat-type')||'tipo',
			active:null,
			items:[]
		};
		loadChatConversations(target);
	});
}

function loadChatConversations(target){
	var c=CHATS[target]; if(!c) return;
	fetch('/api/'+target).then(function(r){return r.json();}).then(function(items){
		c.items=items||[];
		var box=$('chat-conv-'+target); if(!box) return;
		box.innerHTML='';
		if(!items||!items.length){ box.innerHTML='<div class="empty-state"><p>Sem conversas</p></div>'; return; }
		items.forEach(function(item){
			var name=item.nome||item.titulo||item.numero||item.id;
			var last=item.ultima_mensagem||item.last_message||'';
			var unread=item.nao_lidas||0;
			var initials=String(name).trim().split(/\s+/).slice(0,2).map(function(part){return part.charAt(0).toUpperCase();}).join('');
			var div=document.createElement('div');
			div.className='chat-conv'+(String(c.active)===String(item.id)?' active':'');
			div.innerHTML='<div class="chat-conv-avatar">'+esc(initials||'#')+'</div><div class="chat-conv-body"><div class="chat-conv-name">'+esc(name)+'</div><div class="chat-conv-last">'+esc(last||'Sem mensagens')+'</div><div class="chat-conv-meta"><span>'+esc(item.status||'')+'</span><span>'+esc(item.numero||'')+'</span></div></div><div>'+(unread?('<span class="chat-unread">'+unread+'</span>'):'')+'</div>';
			div.onclick=function(){openChat(target,item.id,name);};
			box.appendChild(div);
		});
		if(!c.active&&items[0]) openChat(target,items[0].id,items[0].nome||items[0].titulo||items[0].numero||items[0].id);
	}).catch(function(){ var box=$('chat-conv-'+target); if(box) box.innerHTML='<div class="empty-state"><p>Erro ao carregar conversas</p></div>'; });
}

function refreshChat(target){
	loadChatConversations(target);
	var c=CHATS[target];
	if(c&&c.active){ loadChatMessages(target,c.active); }
}

function openChat(target,id,title){
	var c=CHATS[target]; if(!c) return;
	c.active=id;
	$('chat-title-'+target).textContent=title||('Conversa #'+id);
	loadChatConversations(target);
	loadChatMessages(target,id);
}

function loadChatMessages(target,id){
	var c=CHATS[target]; if(!c) return;
	fetch('/api/'+c.messages+'?'+encodeURIComponent(c.relation)+'='+encodeURIComponent(id)).then(function(r){return r.json();}).then(function(items){
		var box=$('chat-msg-'+target); if(!box) return;
		box.innerHTML='';
		if(!items||!items.length){ box.innerHTML='<div class="empty-state"><p>Sem mensagens</p></div>'; return; }
		items.forEach(function(item){
			var mine=!!item[c.authorField];
			var text=item[c.textField]||'';
			var media=item[c.mediaField]||'';
			var type=(item[c.typeField]||'').toLowerCase();
			var div=document.createElement('div');
			div.className='chat-bubble '+(mine?'mine':'other');
			var html='';
			if(text) html+='<div>'+esc(text)+'</div>';
			if(media){ html+='<div class="chat-media">'+chatMediaHTML(media,type)+'</div>'; }
			html+='<div class="chat-bubble-meta">'+esc(item[c.timeField]||'')+'</div>';
			div.innerHTML=html;
			box.appendChild(div);
		});
		box.scrollTop=box.scrollHeight;
	});
}

function chatMediaHTML(path,type){
	var src='/media/stream?path='+encodeURIComponent(path);
	if((type||'').indexOf('audio')>=0||String(path).match(/\.(mp3|wav|ogg|m4a|aac)$/i)) return '<audio controls src="'+src+'"></audio>';
	if((type||'').indexOf('video')>=0||String(path).match(/\.(mp4|webm|mov|avi)$/i)) return '<video controls src="'+src+'"></video>';
	if(String(path).match(/\.(jpg|jpeg|png|gif|webp|svg)$/i)) return '<img src="'+esc(path)+'" alt="midia">';
	return '<a class="link" target="_blank" href="'+esc(path)+'">Abrir arquivo</a>';
}

function chatSend(target,e){
	e.preventDefault();
	var c=CHATS[target]; if(!c||!c.active) return;
	var input=$('chat-input-'+target); var text=input.value.trim();
	var payload={};
	payload[c.relation]=c.active;
	payload[c.textField]=text;
	payload[c.authorField]=true;
	payload[c.typeField]='chat';
	fetch('/api/'+c.messages,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(payload)})
		.then(function(r){if(!r.ok) throw new Error('Erro ao enviar'); return r.json();})
		.then(function(){ input.value=''; chatTyping(target,false); loadChatMessages(target,c.active); loadChatConversations(target); })
		.catch(function(err){ toast(err.message,'erro'); });
}

function chatUpload(target,input){
	var c=CHATS[target]; if(!c||!c.active||!input.files||!input.files[0]) return;
	var fd=new FormData(); fd.append('file', input.files[0]);
	fetch('/upload',{method:'POST',body:fd}).then(function(r){if(!r.ok) throw new Error('Upload falhou'); return r.json();}).then(function(d){
		var payload={};
		payload[c.relation]=c.active;
		payload[c.textField]=input.files[0].name;
		payload[c.mediaField]=d.path;
		payload[c.authorField]=true;
		payload[c.typeField]=input.files[0].type||'arquivo';
		return fetch('/api/'+c.messages,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(payload)});
	}).then(function(){ input.value=''; loadChatMessages(target,c.active); loadChatConversations(target); }).catch(function(err){ toast(err.message,'erro'); });
}

function chatFilter(target,q){
	q=(q||'').toLowerCase();
	document.querySelectorAll('#chat-conv-'+target+' .chat-conv').forEach(function(el){
		el.style.display=el.textContent.toLowerCase().includes(q)?'':'none';
	});
}

function chatTyping(target,typing){
	var c=CHATS[target]; if(!c||!c.active) return;
	fetch('/api/_presence',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({user:'local',screen:target,ticket_id:c.active,typing:typing,status:typing?'digitando':'online'})}).catch(function(){});
}

function chatHandlePresence(msg){
	var data=msg.data||{};
	var ticketId=String(data.ticket_id||'');
	Object.keys(CHATS).forEach(function(target){
		var c=CHATS[target];
		if(!c||String(c.active)!==ticketId) return;
		if(msg.type==='digitando'){
			$('chat-typing-'+target).textContent=(data.user||'Alguem')+' está digitando...';
			setTimeout(function(){ if($('chat-typing-'+target).textContent.indexOf('digitando')>=0) $('chat-typing-'+target).textContent=''; }, 1800);
		}else if(msg.type==='presenca'){
			$('chat-presence-'+target).textContent=(data.user||'')+' '+(data.status||'');
		}else if(msg.type==='presenca_socket'){
			$('chat-presence-'+target).textContent='Conexões ativas: '+((data&&data.connections)||0);
		}
	});
}

function chatHandleUpdate(msg){
	Object.keys(CHATS).forEach(function(target){
		var c=CHATS[target];
		if(!c) return;
		if(msg.model===c.target){ loadChatConversations(target); }
		if(msg.model===c.messages&&c.active){ loadChatMessages(target,c.active); }
	});
}

// ===== Init =====
document.addEventListener('DOMContentLoaded',function(){
  connectWS();
  renderCharts();
	initChats();
	carregarListasInline();
	refreshSystemStatus();
	setInterval(refreshSystemStatus,8000);
`)
	for _, model := range s.Program.Models {
		b.WriteString(fmt.Sprintf("  carregar('%s');\n", lo(model.Name)))
	}
	b.WriteString("});\n")

	// Auth JS
	if s.Auth != nil {
		loginField := "email"
		passField := "senha"
		if s.Program.Auth != nil {
			if s.Program.Auth.LoginField != "" {
				loginField = s.Program.Auth.LoginField
			}
			if s.Program.Auth.PassField != "" {
				passField = s.Program.Auth.PassField
			}
		}
		b.WriteString(fmt.Sprintf(`
// ===== Auth =====
var authToken=localStorage.getItem('flang_token')||'';
var authMode='login';
var AUTH_LOGIN='%s',AUTH_PASS='%s';

function authHeaders(){
  var h={'Content-Type':'application/json'};
  if(authToken)h['Authorization']='Bearer '+authToken;
  return h;
}

// Override fetch to inject auth token
var _fetch=window.fetch;
window.fetch=function(url,opts){
  opts=opts||{};
  if(authToken&&url.startsWith('/api/')){
    opts.headers=opts.headers||{};
    if(typeof opts.headers==='object'&&!opts.headers['Authorization']){
      opts.headers['Authorization']='Bearer '+authToken;
    }
  }
  return _fetch(url,opts);
};

function mostrarLogin(){
  $('auth-modal').style.display='flex';
  authMode='login';
  $('auth-title').textContent='Entrar';
  $('auth-toggle-text').textContent='Nao tem conta?';
  $('auth-toggle-link').textContent='Criar conta';
  $('auth-error').style.display='none';
  $('auth-extra-fields').innerHTML='';
}
function fecharAuth(){$('auth-modal').style.display='none';}
function toggleAuthMode(){
  if(authMode==='login'){
    authMode='register';
    $('auth-title').textContent='Criar Conta';
    $('auth-toggle-text').textContent='Ja tem conta?';
    $('auth-toggle-link').textContent='Entrar';
    $('auth-extra-fields').innerHTML='<div class="form-group"><label>Nome</label><input type="text" id="auth-nome" required class="form-input" placeholder="Seu nome"></div>';
  } else {
    authMode='login';
    $('auth-title').textContent='Entrar';
    $('auth-toggle-text').textContent='Nao tem conta?';
    $('auth-toggle-link').textContent='Criar conta';
    $('auth-extra-fields').innerHTML='';
  }
  $('auth-error').style.display='none';
}
function authSubmit(e){
  e.preventDefault();
  var login=$('auth-login').value;
  var pass=$('auth-pass').value;
  var url=authMode==='login'?'/api/login':'/api/registro';
  var body={};
  body[AUTH_LOGIN]=login;
  body[AUTH_PASS]=pass;
  if(authMode==='register'){
    var nome=$('auth-nome');
    if(nome)body['nome']=nome.value;
  }
  _fetch(url,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)})
    .then(function(r){return r.json().then(function(d){return{ok:r.ok,data:d};});})
    .then(function(res){
      if(!res.ok){
        $('auth-error').textContent=res.data.erro||'Erro';
        $('auth-error').style.display='block';
        return;
      }
      authToken=res.data.token;
      localStorage.setItem('flang_token',authToken);
      fecharAuth();
      updateAuthUI(res.data.login||login,res.data.role||'usuario');
      toast(authMode==='login'?'Bem-vindo!':'Conta criada!');
      // Reload all data with auth
      Object.keys(M).forEach(function(m){carregar(m);});
      carregarListasInline();
      renderCharts();
    });
}
function sair(){
  authToken='';
  localStorage.removeItem('flang_token');
  $('btn-login').style.display='';
  $('user-info').style.display='none';
  $('btn-logout').style.display='none';
  toast('Desconectado');
}
function updateAuthUI(login,role){
  $('btn-login').style.display='none';
  $('user-info').style.display='inline';
  $('user-info').textContent=login+' ('+role+')';
  $('btn-logout').style.display='inline';
}
// Check stored token on load
if(authToken){
  _fetch('/api/me',{headers:{'Authorization':'Bearer '+authToken}})
    .then(function(r){return r.json();})
    .then(function(d){if(d&&d.login)updateAuthUI(d.login,d.role||'usuario');else sair();})
    .catch(function(){sair();});
}
`, loginField, passField))
	}

	return b.String()
}

// ============================================================
// Helpers
// ============================================================

// fieldTypeCode returns a short code for JS metadata.
func fieldTypeCode(f *ast.Field) string {
	switch f.Type {
	case ast.FieldNumero:
		return "n"
	case ast.FieldDinheiro:
		return "d"
	case ast.FieldStatus:
		return "s"
	case ast.FieldEmail:
		return "e"
	case ast.FieldTextoLongo:
		return "tl"
	case ast.FieldImagem, ast.FieldUpload, ast.FieldArquivo:
		return "f"
	case ast.FieldEnum:
		return "en"
	case ast.FieldBooleano:
		return "b"
	case ast.FieldSenha:
		return "pw"
	default:
		return "t"
	}
}

// enumJSArray builds a JS array literal from enum values.
func enumJSArray(vals []string) string {
	var parts []string
	for _, v := range vals {
		parts = append(parts, fmt.Sprintf("'%s'", v))
	}
	return "[" + strings.Join(parts, ",") + "]"
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
