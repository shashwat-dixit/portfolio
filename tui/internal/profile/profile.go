// Package profile holds the static portfolio content shown in the TUI.
// Keep this in sync with web/src/data/resume.tsx.
package profile

const (
	Name        = "Shashwat Dixit"
	Initials    = "SD"
	URL         = "https://shashwatdixit.com"
	Location    = "Bengaluru, India"
	Email       = "shashwatmain@gmail.com"
	SSHCommand  = "ssh shashwatdixit.com -p 2222"
	Description = "Software Engineer building performant backends and full-stack applications. I care about systems that scale and developer experience that doesn't suck."
	Summary     = "I'm a software engineer at [Interview Kickstart](https://interviewkickstart.com) working on payments — gateway integrations, installment billing, and the flow that grants learners access after they pay. Previously at [Instahyre](https://instahyre.com) I built distributed systems, search infrastructure, and optimized backend performance. Before that I built event-driven pipelines and SSR frontends as a [full-stack contractor](https://shashwatdixit.com/#work). I hold a degree in [Electrical & Electronics Engineering from NMIT Bengaluru](https://nmit.ac.in) and have authored [IEEE research papers on Quantum Computing](https://scholar.google.com/citations?user=q3MbjLQAAAAJ&hl=en). I like working across the stack — from Redis locking to Elasticsearch query parsers to Next.js SPAs."
	CalURL      = "https://cal.com/shashwatdixit/meeting"
)

type Job struct {
	Company  string
	URL      string
	Location string
	Title    string
	Start    string
	End      string
	Badge    string
	Bullets  []string
}

type School struct {
	Name   string
	URL    string
	Degree string
	Start  string
	End    string
}

type Project struct {
	Title        string
	URL          string
	Dates        string
	Description  string
	Technologies []string
	SourceURL    string
}

type Link struct {
	Name string
	URL  string
}

var Work = []Job{
	{
		Company:  "Interview Kickstart",
		URL:      "https://interviewkickstart.com",
		Location: "Bengaluru, India",
		Title:    "SDE-1",
		Start:    "June 2026",
		End:      "Present",
		Bullets: []string{
			"Fixed 5 payment-gateway defects against the live API format, stopping real transactions from silently falling back to a backup provider",
			"Fixed a double-counting bug in installment balance logic that was silently rejecting valid payments",
			"Shipped a feature-flagged payment-flow rewrite with old and new systems running side-by-side for instant rollback",
			"Cut learner activation delay by granting access on payment completion, with a 30-day auto-revoke if payment is never confirmed",
			"Cut error-monitoring noise by ~30% by filtering expected no-ops, and added reference-ID logging for faster debugging",
		},
	},
	{
		Company:  "Instahyre",
		URL:      "https://instahyre.com",
		Location: "Bengaluru, India",
		Title:    "SDE-1",
		Start:    "December 2024",
		End:      "May 2026",
		Bullets: []string{
			"Designed Redis-based distributed locking (SETNX + TTL) to eliminate double-booking under concurrent traffic",
			"Reduced API latency by 80% (p75: 600ms → 120ms) by profiling slow PostgreSQL queries and optimizing payloads",
			"Built a stack-based boolean query parser (AND/OR/NOT) for Elasticsearch, improving matching relevance by 20%",
			"Reduced regression bugs by 30% by decomposing a monolith into independently maintainable modules",
		},
	},
	{
		Company:  "Pummyz Foods",
		URL:      "https://upwork.com",
		Location: "Remote",
		Title:    "Full-Stack Developer",
		Start:    "May 2024",
		End:      "August 2024",
		Badge:    "Contract",
		Bullets: []string{
			"Designed an event-driven pipeline using Kafka for asynchronous order processing across thousands of daily events",
			"Implemented idempotent webhook-based payment confirmation, reducing duplicate order processing by 95%",
			"Reduced media load times by 30% via AWS S3 + CDN caching",
			"Improved Core Web Vitals by 35% with a Next.js SSR architecture optimized for first-contentful paint and reduced JS bundle size",
		},
	},
}

var Education = []School{
	{
		Name:   "Nitte Meenakshi Institute of Technology, Bengaluru",
		URL:    "https://nmit.ac.in",
		Degree: "Bachelor of Engineering in Electrical & Electronics Engineering",
		Start:  "2020",
		End:    "2024",
	},
}

var Skills = []string{
	"Python", "C++", "Go", "JavaScript", "TypeScript", "Git", "GCP", "BigQuery",
	"AWS", "Docker", "Kubernetes", "React", "Next.js", "Node.js", "Bun.js",
	"Elasticsearch", "Django", "RabbitMQ", "Kafka", "PostgreSQL", "Redis",
	"Cassandra", "MySQL", "SQL",
}

var Projects = []Project{
	{
		Title:       "Jamin",
		URL:         "https://jamin.shashwatdixit.com",
		Dates:       "2024",
		Description: "Full-stack AI chat platform integrating multiple LLMs (GPT, Claude, Cohere) via LangChain. Features RAG-powered PDF Q&A and YouTube summarization, plus Stable Diffusion image generation. Scalable PostgreSQL/Drizzle backend enabling context-aware conversations across models.",
		Technologies: []string{
			"React", "Node.js", "LangChain", "PostgreSQL", "Drizzle", "OpenAI", "Stable Diffusion",
		},
		SourceURL: "https://github.com/shashwat-dixit/jamin",
	},
	{
		Title:       "Zort",
		URL:         "https://zort.shashwatdixit.com",
		Dates:       "2024",
		Description: "Real-time collaborative whiteboard supporting 100+ concurrent users, built with Next.js and Socket.IO. Client-side state management via Zustand and localStorage to reduce server load. CI/CD pipeline with Docker and GitHub Actions for zero-failure automated deployments.",
		Technologies: []string{
			"Next.js", "Socket.IO", "Zustand", "Docker", "GitHub Actions",
		},
		SourceURL: "https://github.com/shashwat-dixit/zort",
	},
	{
		Title:       "Phabric",
		URL:         "https://phabric.shashwatdixit.com",
		Dates:       "2024",
		Description: "Phabric is Vercel but with websocket support!",
		Technologies: []string{
			"Next.js", "Socket.IO", "Zustand", "Docker", "GitHub Actions",
		},
		SourceURL: "https://github.com/shashwat-dixit/zort",
	},
	{
		Title:       "Code Compete",
		URL:         "https://codecompete.shashwatdixit.com",
		Dates:       "2024",
		Description: "Code Compete is a platform for competitive programming problems. It is built with Next.js and Tailwind CSS.",
		Technologies: []string{
			"Next.js", "Socket.IO", "Zustand", "Docker", "GitHub Actions",
		},
		SourceURL: "https://github.com/shashwat-dixit/zort",
	},
}

var Social = []Link{
	{Name: "GitHub", URL: "https://github.com/shashwat-dixit"},
	{Name: "LinkedIn", URL: "https://linkedin.com/in/dixitshashwat"},
	{Name: "X", URL: "https://x.com/shashwatmain"},
	{Name: "Email", URL: "mailto:" + Email},
	{Name: "Schedule a meeting", URL: CalURL},
}
