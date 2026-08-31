name: Comment
fields:
  ID:
    type: AutoIncrement
  Text:
    type: String
identifiers:
  primary: ID
related:
  Commentable:
    type: ForOnePoly
    for:
      - Person
      - Company
