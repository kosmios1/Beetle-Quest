#let report(
  title: "",
  subtitle: none,
  authors: ("",),
  date: none,
  doc,
  imagepath: "",
) = {
  set document(title: title, author: authors)

  set page(
    numbering: "1",
    paper: "a4",
    margin: 6em,
  )
  set text(size: 12.5pt, font: "New Computer Modern")
  show raw: set text(size: 12pt, font: "IosevkaTerm NF")
  
  set par(
    justify: true,
    leading: 0.55em,
    linebreaks: "optimized",
    first-line-indent: 1.8em,
  )
  show par: set block(spacing: 0.55em)

  set heading(numbering: "1.")
  show heading: set block(above: 1.4em, below: 1em)
  
  set math.equation(numbering: "(1)")

  align(center)[
    #image(imagepath, width: 25em)
  ]

  align(horizon + center)[
    #text(25pt, title, weight: "bold")
    #v(1em)
    #text([#subtitle], weight: "regular")
    #v(0.5em)
    #text([#date], weight: "regular")
    #v(2em)
    #text(
      list(..authors, marker: "", body-indent: 0pt), style:"italic" )
  ]

  // Table of contents
  align(bottom)[
    #show outline.entry.where(level: 1): it => { strong(it) }
    #outline(title: [Index #v(1em)], indent: 2em)
  ]

  pagebreak()

  doc
}
