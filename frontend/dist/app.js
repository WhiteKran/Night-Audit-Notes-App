const addBtn = document.getElementById('addBtn');
const notesList = document.getElementById('notesList');
const lockedNotes = new Set();

// Prevent zoom
document.addEventListener('keydown', (e) => {
	if ((e.ctrlKey || e.metaKey) && (e.key === '+' || e.key === '-' || e.key === '=' || e.key === '0')) {
		e.preventDefault();
	}
});

document.addEventListener('wheel', (e) => {
	if (e.ctrlKey || e.metaKey) {
		e.preventDefault();
	}
}, { passive: false });

window.addEventListener('wails:ready', async () => {
	await loadNotes();
});

addBtn.addEventListener('click', async () => {
	await window.go.main.App.AddNote('');
	await loadNotes();
});

document.addEventListener("DOMContentLoaded", () => {
	loadNotes();
});

async function loadNotes() {
	try {
		const notes = await window.go.main.App.GetNotes();
		renderNotes(notes);
	} catch (err) {
		console.error('Error loading notes:', err);
	}
}

function formatDate(dateString) {
	const date = new Date(dateString);
	return date.toLocaleString();
}

function renderNotes(notes) {
	notesList.innerHTML = notes.map((note) => {
		const updatedAt = formatDate(note.updatedAt);
		const isLocked = lockedNotes.has(note.id);

		return `
			<div class="note-item" data-note-id="${note.id}">
				<div class="note-actions">
					<button class="btn btn-copy" data-note-id="${note.id}"><img src="res/copy.svg" alt="Copy" class="icon"/></button>
				</div>
				<div class="note-content">
					<div class="textarea-border">
						<img src="res/corner.svg" class="corner corner-tl" alt=""/>
						<img src="res/corner.svg" class="corner corner-tr" alt=""/>
						<img src="res/corner.svg" class="corner corner-bl" alt=""/>
						<img src="res/corner.svg" class="corner corner-br" alt=""/>
						<img src="res/side-top.svg" class="side side-top" alt=""/>
						<img src="res/side-bottom.svg" class="side side-bottom" alt=""/>
						<img src="res/side-left.svg" class="side side-left" alt=""/>
						<img src="res/side-right.svg" class="side side-right" alt=""/>
						<textarea class="note-text-editable" data-note-id="${note.id}" data-original-text="${escapeHtml(note.text)}">${note.text}</textarea>
					</div>
				</div>
				<div class="note-actions">
					<button class="btn btn-lock ${isLocked ? 'locked' : ''}" data-note-id="${note.id}"><img src="res/${isLocked ? 'locked' : 'unlocked'}.svg" alt="Lock" class="icon"/></button>
					<button class="btn btn-remove" data-note-id="${note.id}" ${isLocked ? 'disabled' : ''}><img src="res/delete.svg" alt="Delete" class="icon"/></button>
				</div>
			</div>
		`;
	}).join('');

	// Attach event listeners
	document.querySelectorAll('.note-text-editable').forEach(textarea => {
		const noteId = parseInt(textarea.dataset.noteId);

		// Auto-save on input
		textarea.addEventListener('input', async (e) => {
			const text = textarea.value;
			try {
				await window.go.main.App.UpdateNote(noteId, text);
			} catch (err) {
				console.error('Failed to save note:', err);
			}
		});

		// Auto-resize textarea
		textarea.addEventListener('input', () => {
			textarea.style.height = 'auto';
			textarea.style.height = Math.min(textarea.scrollHeight, 400) + 'px';
		});

		// Initial resize
		setTimeout(() => {
			textarea.style.height = 'auto';
			textarea.style.height = Math.min(textarea.scrollHeight, 400) + 'px';
		}, 0);
	});

	document.querySelectorAll('.btn-copy').forEach(btn => {
		const icon = btn.querySelector('img.icon');

		btn.addEventListener('click', async (e) => {
			const noteId = parseInt(btn.dataset.noteId);
			const notes = await window.go.main.App.GetNotes();
			const note = notes.find(n => n.id === noteId);
			if (note) {
				try {
					await navigator.clipboard.writeText(note.text);
					icon.style.filter = 'invert(64%) sepia(98%) saturate(451%) hue-rotate(76deg) brightness(91%) contrast(92%)';
					setTimeout(() => {
						icon.style.filter = '';
					}, 150);
				} catch (err) {
					console.error('Failed to copy:', err);
				}
			}
		});
	});

	document.querySelectorAll('.btn-remove').forEach(btn => {
		const icon = btn.querySelector('img.icon');

		btn.addEventListener('click', async (e) => {
			const noteId = parseInt(btn.dataset.noteId);
			const notes = await window.go.main.App.GetNotes();
			const index = notes.findIndex(n => n.id === noteId);
			const note = notes[index];
			if (note) {
				deletedNote = { note, index };
			}
			try {
				await window.go.main.App.RemoveNote(noteId);
				await loadNotes();
			} catch (err) {
				console.error('Failed to remove note:', err);
				deletedNote = null;
			}
		});

		btn.addEventListener('mouseenter', () => {
			icon.src = 'res/delete_hover.svg';
		});

		btn.addEventListener('mouseleave', () => {
			icon.src = 'res/delete.svg';
		});
	});

	document.querySelectorAll('.btn-lock').forEach(btn => {
		btn.addEventListener('click', async (e) => {
			const noteId = parseInt(btn.dataset.noteId);
			const notes = await window.go.main.App.GetNotes();

			if (lockedNotes.has(noteId)) {
				lockedNotes.delete(noteId);
			} else {
				lockedNotes.add(noteId);
			}

			renderNotes(notes);
		});
	});
}

function escapeHtml(text) {
	const map = {
		'&': '&amp;',
		'<': '&lt;',
		'>': '&gt;',
		'"': '&quot;',
		"'": '&#039;'
	};
	return text.replace(/[&<>"']/g, m => map[m]);
}
