# HARFANG® Tutorials

These **tutorials** demonstrate the usage of the HARFANG API in **Python** and **Lua**.

To run the tutorials:

1. Download or clone this repository to your computer _(eg. in `d:/tutorials-hg2`)_.
2. Download _assetc_ for your platform from [here](https://harfang3d.com/releases) to compile the tutorial resources.
3. Drag and drop the tutorial resources folder on the assetc executable **-OR-** execute assetc passing it the path to the tutorial resources folder _(eg. `assetc d:/tutorials-hg2/resources`)_.

![assetc drag & drop](https://github.com/harfang3d/image-storage/raw/main/tutorials/assetc.gif)

After the compilation process finishes, you will see a `resources_compiled` folder next to the tutorials resources folder.

You can now execute the tutorials from the folder you unzipped them to.

```bash
D:\tutorials-hg2>python draw_lines.py
```
or
```bash
D:\tutorials-hg2>lua draw_lines.lua
```

Alternatively, you can open the tutorial folder using [Visual Studio Code](https://code.visualstudio.com/) and use the provided debug targets to run the tutorials.

**If you want to know more about HARFANG**, please visit the [official website](https://www.harfang3d.com).

## Screenshots
* AAA Rendering Pipeline
![AAA rendering pipeline](screenshots/aaa.png)

* Mouse flight controller
![Mouse flight controller](screenshots/game_mouse_flight.png)

* PBR materials
![PBR materials](screenshots/scene_pbr.png)

* Draw to multiple viewports
![Physics pool](screenshots/scene_draw_to_multiple_viewports.png)

* Physics pool of objects
![Physics pool](screenshots/physics_pool_of_objects.png)

* Physics Kapla
![Physics Kapla](screenshots/physics_kapla.png)

* Scene Instances
![Scene Instances](screenshots/scene_instances.png)

* Scene with many nodes
![Many object nodes](screenshots/scene_many_nodes.png)

* Forward rendering pipeline, using the lights priority
![Lights priority](screenshots/scene_light_priority.png)

* Dear ImGui edition
![Dear ImGui edition](screenshots/imgui_edit.png)
