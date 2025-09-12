# Gradebook API Naming Conventions

This document outlines the consistent naming conventions used throughout the gradebook system for both JSON APIs and internal data structures.

## Core Principles

1. **snake_case for JSON fields** - All JSON field names use lowercase with underscores
2. **Semantic consistency** - JSON field names align with their content semantics  
3. **No stutter** - Avoid repetitive naming like `grade.Grade` in code access patterns
4. **Descriptive struct names** - Use clear, specific names that indicate purpose

## JSON Field Naming

### Assignment-level Fields
All assignment metadata uses the `assignment_` prefix:
- `assignment_name` - Human-readable assignment name
- `assignment_date` - Date in YYYYMMDD format  
- `assignment_type` - Type of assignment (quiz, test, hw, etc.)
- `assignment_category` - Grading category (major, minor, cp, etc.)
- `assignment_records` - Array of student grade records

### Student-level Fields
Student information uses descriptive field names:
- `first_name` - Student's first name
- `last_name` - Student's last name
- `email` - Student's email address (used as unique identifier)

### Grade Records
Individual grade records within `assignment_records`:
- `email` - Student identifier
- `grade` - Numeric grade value (can be null for absent students)

### Class Configuration Fields
Class-level mappings use the pattern `{content}_by_{key}`:
- `students_by_email` - Map of email → student info
- `categories_by_assignment_type` - Map of assignment type → category
- `labels_by_assignment_category` - Map of category → human-readable label
- `weights_by_assignment_category` - Map of category → numeric weight
- `terms_by_id` - Map of term ID → term date range

## Data Structure Naming (Go)

### Primary Types
- `AssignmentRecord` - Individual student grade record
- `AssignmentRecords` - Array of AssignmentRecord structs  
- `Gradebook` - Single assignment's complete data
- `Class` - Complete class configuration and student roster
- `Student` - Individual student information and calculated grades

### Field Mapping
Go struct fields use PascalCase but map to snake_case JSON:

```go
type AssignmentRecord struct {
    Email string   `json:"email"`
    Grade *float64 `json:"grade"`
}

type Gradebook struct {
    AssignmentName     string            `json:"assignment_name"`
    AssignmentDate     string            `json:"assignment_date"`
    AssignmentType     string            `json:"assignment_type"`
    AssignmentCategory string            `json:"assignment_category"`
    AssignmentRecords  AssignmentRecords `json:"assignment_records"`
}
```

## File Structure

### Class Configuration (`class.json`)
Contains class-wide settings, student roster, grading categories, and term definitions.

### Assignment Data (`*.gradebook`)  
Individual assignment files containing:
- Assignment metadata (name, date, type, category)
- Array of student grade records
- Each record links student (by email) to grade value

## Consistency Benefits

1. **Predictable patterns** - Developers can infer field names
2. **No semantic conflicts** - Grade records contain "grade" fields
3. **Cross-language portability** - JSON naming works across implementations
4. **Maintainable codebase** - Clear relationship between JSON and internal structures

## Migration Notes

When adapting this system to other languages:
- Maintain snake_case JSON field names exactly
- Adapt struct/class names to language conventions (e.g., Python: assignment_record)
- Preserve the semantic relationships between collections and items
- Keep the `assignment_records` → record with `grade` field consistency
