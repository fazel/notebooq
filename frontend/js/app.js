let token = localStorage.getItem("token");

// Redirect to notes if already logged in
if (token && window.location.pathname.endsWith("login.html")) {
    window.location.href = "/notes.html";
}

// Login
const loginForm = document.getElementById("loginForm");
if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const username = document.getElementById("loginUsername").value;
        const password = document.getElementById("loginPassword").value;
        try {
            const res = await axios.post("/api/login", { username, password });
            token = res.data.token;
            localStorage.setItem("token", token);
            window.location.href = "/notes.html";
        } catch (err) {
            alert(err.response?.data?.error || err.message);
        }
    });
}

// Signup
const signupForm = document.getElementById("signupForm");
if (signupForm) {
    signupForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const username = document.getElementById("signupUsername").value;
        const password = document.getElementById("signupPassword").value;
        try {
            await axios.post("/api/signup", { username, password });
            alert("Signup successful! You can login now.");
            window.location.href = "/login.html";
        } catch (err) {
            alert(err.response?.data?.error || err.message);
        }
    });
}

// Notes page
const notesList = document.getElementById("notesList");
const noteForm = document.getElementById("noteForm");
const logoutBtn = document.getElementById("logoutBtn");

if (token) {
    axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;

    async function loadNotes() {
        try {
            const res = await axios.get("/api/notes");
            notesList.innerHTML = res.data.map(n => `
                <div class="card mb-2">
                    <div class="card-body">
                        <h5 class="card-title">${n.title}</h5>
                        <p class="card-text">${n.content}</p>
                    </div>
                </div>
            `).join("");
        } catch (err) {
            console.error(err);
        }
    }

    loadNotes();

    if (noteForm) {
        noteForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            const title = document.getElementById("noteTitle").value;
            const content = document.getElementById("noteContent").value;
            try {
                await axios.post("/api/notes", { title, content });
                noteForm.reset();
                loadNotes();
            } catch (err) {
                alert(err.response?.data?.error || err.message);
            }
        });
    }

    if (logoutBtn) {
        logoutBtn.addEventListener("click", () => {
            localStorage.removeItem("token");
            window.location.href = "/login.html";
        });
    }
}
