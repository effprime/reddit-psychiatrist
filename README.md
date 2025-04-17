# 🧠 reddit-psychiatrist

**Analyze any Reddit user’s personality based on their public comments.**  
It’s part psychoanalysis, part roast, part mirror you didn’t ask for.

Built with Go, powered by GPT-4, and brutally honest.

---

## 🚨 What It Does

- 🔍 Scrapes a Reddit user’s public comments
- 🧠 Uses GPT-4 to infer:
  - Personality summary (no fluff, mostly roast)
  - Core interests (inferred, not just subreddits)
- 💀 Returns an unfiltered psychological profile
- 🌐 Includes both:
  - CLI interface
  - Web interface (basic Go templates)
  - API mode (for automation or frontend integration)

---

## 🛠️ Tech Stack

- **Go** — backend and CLI
- **OpenAI GPT-4o** — personality engine
- **Reddit public API** — no auth required
- **HTML templates** — zero-JS UI for now
- **net/http** — no framework, just vibes

---

## 📦 Installation

```bash
git clone https://github.com/yourusername/reddit-psychiatrist.git
cd reddit-psychiatrist
go build ./cmd/...
