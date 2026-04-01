# `pptx-probe`

`pptx-probe` is a simple automation tool that transforms **Microsoft PowerPoint** into a surprisingly capable IDE (I know it sounds stupid, but it's indeed stupid.) for technical presentations and rapid prototyping. It is designed to bridge the gap between static slides and active codebases by treating slide text boxes as source files.



### Why use PowerPoint as an IDE?
In academic settings (like **NUS CS2030S** labs), code is often shown on slides. `pptx-probe` ensures that the code your audience sees is the exact same code being compiled and benchmarked in the background, eliminating "copy-paste" errors and desynced documentation.

### Checklist

- [x] Add `watch` command to automatically update the code when saving the powerpoint file.
- [ ] Perform syntax hightlighting when saved.

### Quick start

```bash
make
./pptx-probe -o output example/example1.pptx 
```

For more details, please refer to [this guide](./Guide.md).

### Example
| | |
|:---:|:---:|
| ![alt text](figure/example_pic_1.png) | ![alt text](figure/example_pic_2.png) |




```cpp
// Under cat.cpp
#include <iostream>
#include "cat.hpp" // Crucial: Link the implementation to the header

// Define the logic for the 'speak' method
void Cat::speak() const {
    std::cout << name << " says: Meow! (I am " << age << " years old)" << std::endl;
}

// Under cat.hpp
#ifndef CAT_HPP
#define CAT_HPP

#include <string>

// The Interface: Tells the compiler what a 'Cat' is.
struct Cat {
    std::string name;
    int age;

    // Prototype: We promise this function exists somewhere else.
    void speak() const;
};

#endif

// Under main.cpp
#include "cat.hpp"

int main() {
    // Create an instance using the definition from the header
    Cat myCat = {"Luna", 3};
    
    // Call the function defined in the implementation file
    myCat.speak();

    return 0;
}


```
