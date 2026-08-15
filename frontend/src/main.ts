import { Events } from "@wailsio/runtime";
import { FastNoteService } from "../bindings/fastnote";

let currentText = "";
let saveTimer: ReturnType<typeof setTimeout>;

// Get initial state from Go backend
const state = await FastNoteService.GetState();
if (state) {
  const editor = document.getElementById("editor") as HTMLTextAreaElement;
  const status = document.getElementById("status");
  const preview = document.getElementById("preview");

  if (state.text) {
    editor.value = state.text;
    currentText = state.text;
  }
  if (state.preview) {
    preview.textContent = state.preview;
  }
  if (state.status) {
    status.textContent = state.status;
  }
  if (state.themeIndex === 1) {
    document.documentElement.setAttribute("data-theme", "dark");
  }
}

// Wire toolbar
document.getElementById("btn-open")?.addEventListener("click", async () => {
  const path = prompt("Enter file path:");
  if (!path) return;
  const result = await FastNoteService.OpenFile(path);
  if (result.error) {
    setStatus(result.error);
    return;
  }
  const editor = document.getElementById("editor") as HTMLTextAreaElement;
  editor.value = result.text;
  currentText = result.text;
  setPreview(result.preview);
  setStatus(result.status);
});

document.getElementById("btn-save")?.addEventListener("click", async () => {
  const result = await FastNoteService.SaveFile();
  setStatus(result.error ? result.error : result.status);
});

document.getElementById("btn-saveas")?.addEventListener("click", async () => {
  const path = prompt("Enter save path:");
  if (!path) return;
  const result = await FastNoteService.SaveAsFile(path);
  setStatus(result.error ? result.error : result.status);
});

document.getElementById("btn-export-html")?.addEventListener("click", async () => {
  const path = prompt("Enter export path:");
  if (!path) return;
  const result = await FastNoteService.ExportHTML(path);
  setStatus(result.error ? result.error : result.status);
});

document.getElementById("btn-export-pdf")?.addEventListener("click", async () => {
  const path = prompt("Enter export path:");
  if (!path) return;
  const result = await FastNoteService.ExportPDF(path);
  setStatus(result.error ? result.error : result.status);
});

document.getElementById("btn-theme")?.addEventListener("click", async () => {
  const result = await FastNoteService.ToggleTheme();
  if (result.theme === "dark") {
    document.documentElement.setAttribute("data-theme", "dark");
  } else {
    document.documentElement.removeAttribute("data-theme");
  }
  setStatus(result.status);
});

// Editor input - debounce sync to Go
const editor = document.getElementById("editor") as HTMLTextAreaElement;
editor?.addEventListener("input", async () => {
  const newText = editor.value;
  if (saveTimer) clearTimeout(saveTimer);
  saveTimer = setTimeout(async () => {
    if (newText !== currentText) {
      currentText = newText;
      const preview = await FastNoteService.SetText(newText);
      setPreview(preview);
    }
  }, 100);
});

function setPreview(text: string) {
  const el = document.getElementById("preview");
  if (el && text) el.textContent = text;
}

function setStatus(text: string) {
  const el = document.getElementById("status");
  if (el) el.textContent = text;
}

// Keyboard accelerators (spec 5.2)
document.addEventListener("keydown", async (e) => {
  const ctrl = e.ctrlKey || e.metaKey;
  const shift = e.shiftKey;

  if (ctrl && !shift && e.key === "o") {
    e.preventDefault();
    document.getElementById("btn-open")?.click();
  } else if (ctrl && !shift && e.key === "s") {
    e.preventDefault();
    document.getElementById("btn-save")?.click();
  } else if (ctrl && shift && e.key === "S") {
    e.preventDefault();
    document.getElementById("btn-saveas")?.click();
  } else if (ctrl && !shift && e.key === "e") {
    e.preventDefault();
    document.getElementById("btn-export-html")?.click();
  } else if (ctrl && shift && e.key === "E") {
    e.preventDefault();
    document.getElementById("btn-export-pdf")?.click();
  }
});
