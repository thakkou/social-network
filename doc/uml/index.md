# UML Diagrams — Social Network

This directory contains UML diagrams for the social network application. All diagrams use [Mermaid.js](https://mermaid.js.org/) syntax and render natively in GitHub, VS Code, and most markdown viewers.

## Diagrams

| # | Diagram | File | Description |
|---|---------|------|-------------|
| 1 | **Use Case Diagram** | [use-case.md](./use-case.md) | Actors (Guest, User, Group Admin) and their interactions with the system |
| 2 | **Class Diagram** | [class.md](./class.md) | Domain models, their attributes, and relationships |
| 3 | **Sequence Diagram** | [sequence.md](./sequence.md) | Interaction flows for: creating a post, sending a message, following a user |
| 4 | **Activity Diagram** | [activity.md](./activity.md) | Workflows for: creating a post, group membership, follow system |
| 5 | **Component Diagram** | [component.md](./component.md) | System architecture: Frontend (Next.js), Backend (Go), Database, WebSocket |
| 6 | **ER Diagram** | [erd.md](./erd.md) | Full database schema with all tables, columns, and relationships |

## How to View

Open any `.md` file in:
- **VS Code** — install the [Mermaid Preview](https://marketplace.visualstudio.com/items?itemName=vstirbu.vscode-mermaid-preview) extension
- **GitHub** — diagrams render automatically in `.md` files
- **Browser** — use the [Mermaid Live Editor](https://mermaid.live/) by pasting the diagram code
