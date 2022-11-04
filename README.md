# goalmap

Roadmap and goal mapping tool in Go.  Takes a nested list in markdown format on stdin, spits out a Graphviz dot file on stdout.

Example input:
```
- eat                                                                                                                         
  - make dinner                                             
    - buy ingredients                                             
      - find recipe                                                 
      - make income                                                    
    - have a place to cook                                             
      - pay electric bill                                             
        - make income                                                                
          - help someone who can pay you                                             
            - learn how to help others                                             
      - buy stove                                                          
        - make income                                                      
      - buy pots and pans                                                  
        - make income                                                      
      - pay rent or buy a houe                                             
        - make income                                 
```

Example output:
<img src="./examples/food.svg">
