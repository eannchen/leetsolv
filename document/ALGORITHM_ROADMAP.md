# Adaptive SM-2 Algorithm


## Version 2: Reasoning & Coding Grades

*(This roadmap is subject to change)*

Currently, the version 1 of Adaptive SM-2 Algorithm uses reasoning factor as a "punisher". If memory is used, the algorithm prolongs the next review schedule, giving users more time to "forget" the problem, and the reasoning factor also reduces the ease factor (score) since the review is less effective.

But users could have a habit of not reasoning before coding, sometimes even though we don't have a direct snapshot of a problem in mind, but we may still notice the "pattern" and code it down directly, especially experienced problem solvers. **We may think it's "reasoned", but actually, it's not.**

To make the reasoning factor more powerful and encourage users to reason every problem before they code, the plan is **to separate "familiarity" into "reasoning grade" and "coding grade"**:
- **Reasoning grade**: how the user can reason the problem correctly
- **Coding grade**: how the user can code the problem with the given reasoning

Both grades will be measured using the existing 5-point scale (`VeryHard` to `VeryEasy`).

With these inputs, **the app will always ask the user: "How well can you correctly reason this problem?"** This way, users can be more self-aware of their reasoning ability and coding ability are actually different skills.

Also, the change provides a more detailed reviewing schedule:

1. It provides a better **ease factor adjustment** formula, when users have a high reasoning grade, a lower coding grade shouldn't lower the ease factor too much. On the other hand, a lower reasoning grade should lower the ease factor more even if the coding grade is high. (e.g., `(0.7 * reasoning_penalty) + (0.3 * coding_penalty)`)

2. It provides a better **Due Priority Scoring** formula, prioritizing problems with a low reasoning grade more than low coding grade. (e.g., `(0.7 * reasoning_grade) + (0.3 * coding_grade)`)

3. Users can search problems by reasoning grade and coding grade. Which is more useful for targeting reviews.

In short, the version 2 change the logic from `mixed_familiarity - memory_punishment` to `reasoning_grade & coding_grade`.