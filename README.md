[English](./README.md) | [繁體中文](./README.zh-TW.md) | [简体中文](./README.zh-CN.md)


# LeetSolv
[![Release](https://img.shields.io/github/release/eannchen/leetsolv.svg)](https://github.com/eannchen/leetsolv/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/eannchen/leetsolv)](https://goreportcard.com/report/github.com/eannchen/leetsolv)
[![CI/CD](https://github.com/eannchen/leetsolv/actions/workflows/ci.yml/badge.svg)](https://github.com/eannchen/leetsolv/actions/workflows/ci.yml)

**LeetSolv** is a CLI tool for **Data Structures and Algorithms (DSA)** problem revision with **spaced repetition**. It supports problems from [LeetCode](https://leetcode.com) and [HackerRank](https://hackerrank.com). Powered by a customized [SuperMemo 2](https://en.wikipedia.org/wiki/SuperMemo) algorithm that incorporates **familiarity**, **importance**, and **reasoning** to move beyond rote memorization.

> ***Zero Dependencies**: Implemented in pure Go with no third-party libraries or external tools—full control over all implementations. See [MOTIVATION.md](document/MOTIVATION.md).*

![Demo](document/image/DEMO_header.gif)

## Table of Contents
- [LeetSolv](#leetsolv)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Scoop (Windows)](#scoop-windows)
    - [Homebrew (macOS/Linux)](#homebrew-macoslinux)
    - [Shell Script (macOS/Linux)](#shell-script-macoslinux)
    - [Verify Installation](#verify-installation)
  - [Review Scheduling System](#review-scheduling-system)
    - [Adaptive SM-2 Algorithm](#adaptive-sm-2-algorithm)
    - [Due Priority Scoring](#due-priority-scoring)
    - [Interval Growing Curve](#interval-growing-curve)
  - [Problem Management](#problem-management)
    - [Functionalities](#functionalities)
    - [Data Privacy \& Safety](#data-privacy--safety)
  - [Usage](#usage)
  - [Configuration](#configuration)
  - [FAQ](#faq)
      - [Q: Why use LeetSolv instead of an Anki deck?](#q-why-use-leetsolv-instead-of-an-anki-deck)
      - [Q: Should I add all my previously solved problems?](#q-should-i-add-all-my-previously-solved-problems)
      - [Q: After a period of use, I accumulated too many due problems.](#q-after-a-period-of-use-i-accumulated-too-many-due-problems)

## Installation

### Scoop (Windows)

```powershell
scoop bucket add eannchen https://github.com/eannchen/scoop-bucket
scoop install leetsolv
```

### Homebrew (macOS/Linux)

```bash
brew tap eannchen/tap
brew install leetsolv
```

### Shell Script (macOS/Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash
```

To uninstall:

```bash
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash -s -- --uninstall
```

### Verify Installation
```bash
leetsolv version
leetsolv help
```

## Review Scheduling System

### Adaptive SM-2 Algorithm

Unlike standard SM-2 (used by Anki), LeetSolv adds **importance** and **reasoning** factors—designed for DSA practice, not flashcard memorization. Familiarity (5 levels), importance (4 levels), and reasoning (3 levels) determine your next review date. Randomization prevents bunching reviews on the same days.

```mermaid
graph TD
    A[Add a DSA Problem to LeetSolv] --> B[Algorithm Applies Adaptations]
    B --> C[Familiarity Scale]
    B --> D[Importance Scale]
    B --> E[Reasoning Scale]
    B --> F[Due Penalty, optional]
    B --> G[Randomization, optional]


    H[Algorithm Calculates with SM-2 Ease Factor]
    C --> H
    D --> H
    E --> H
    F --> H
    G --> H

    H --> I[Determine Next Review]
```

### Due Priority Scoring
Due reviews can accumulate over time. LeetSolv ranks them by priority score so you can focus on what matters most.

> *Default formula: (1.5×Importance) + (0.5×Overdue Days) + (3.0×Familiarity) + (-1.5×Review Count) + (-1.0×Ease Factor)*

![Demo](document/image/DEMO_due_scoring.gif)

### Interval Growing Curve

Review intervals expand based on importance, familiarity, and reasoning. Higher importance = shorter intervals, more frequent reviews.

![SM2 Critical](document/image/SM2_CRITICAL.png)
![SM2 High](document/image/SM2_HIGH.png)
![SM2 Medium](document/image/SM2_MEDIUM.png)
![SM2 Low](document/image/SM2_LOW.png)

## Problem Management

### Functionalities

- **CRUD + Undo**: Create, view, update, delete problems. Undo your last action.
- **Trie-Based Search**: Fast filtering by keyword, importance, familiarity.
- **Quick Views**: Summary of due/upcoming problems with paginated listing.
- **Interactive & Batch Modes**: Run interactively or pass commands directly.
- **Intuitive Commands**: Familiar aliases (`ls`, `rm`), color-coded output.
![Demo](document/image/DEMO_mgmt.gif)

### Data Privacy & Safety

- **No Data Collection**: LeetSolv does not upload user data to the internet.
- **Atomic Writes**: All writes use temp file + rename for data consistency.

## Usage

LeetSolv can be run interactively or by passing commands directly from your terminal.

```bash
# Start interactive mode
leetsolv

# Or run commands directly
leetsolv add https://leetcode.com/problems/two-sum
leetsolv status

# Get help
leetsolv help
```

[View Full Usage Guide (USAGE.md)](document/USAGE.md)

## Configuration

Customize via environment variables or JSON config. See [CONFIGURATION.md](document/CONFIGURATION.md) for all options.

## FAQ

#### Q: Why use LeetSolv instead of an Anki deck?

A: Anki is great for memorizing facts, but DSA requires deeper practice. LeetSolv's SM-2 algorithm uses reasoning, familiarity, and importance to schedule deliberate problem-solving—not rote recall.

#### Q: Should I add all my previously solved problems?

A: No. Only add problems you want to revisit. The algorithm uses the add date for scheduling—bulk-adding creates unrealistic schedules. For old problems, re-solve first, then add.

#### Q: After a period of use, I accumulated too many due problems.

A: SM-2 accumulates dues if you skip days. Use [Due Priority Scoring](#due-priority-scoring) to focus on high-priority problems first. Remove mastered problems—the goal is active practice, not tracking everything.

[Open an issue](https://github.com/eannchen/leetsolv/issues) for questions or suggestions.

---

**LeetSolv** - A spaced repetition CLI for DSA, powered by a custom SM-2 algorithm for deliberate practice.