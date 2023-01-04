-- Create and draw models without a pipeline and using ModelBuilder

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Draw and Create Models using ModelBuilder - no Pipeline',
                    res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

-- vertex layout, materials and models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

mdl_builder = hg.ModelBuilder()

-- Below are the 6 faces of the cube model, with 4 vertices declared for each face, that are added as 2 triangles to make a complete quad face

-- -
vertex0 = hg.Vertex() -- Vertex constructor
vertex0.pos = hg.Vec3(-0.5, -0.5, -0.5) -- Declaring the vertex position
vertex0.normal = hg.Vec3(0, 0, -1) -- Declaring the normal
vertex0.uv0 = hg.Vec2(0, 0) -- Declaring the UV
a = mdl_builder:AddVertex(vertex0) -- Adding the vertex to the ModelBuilder

vertex1 = hg.Vertex()
vertex1.pos = hg.Vec3(-0.5, 0.5, -0.5)
vertex1.normal = hg.Vec3(0, 0, -1)
vertex1.uv0 = hg.Vec2(0, 1)
b = mdl_builder:AddVertex(vertex1)

vertex2 = hg.Vertex()
vertex2.pos = hg.Vec3(0.5, 0.5, -0.5)
vertex2.normal = hg.Vec3(0, 0, -1)
vertex2.uv0 = hg.Vec2(1, 1)
c = mdl_builder:AddVertex(vertex2)

vertex3 = hg.Vertex()
vertex3.pos = hg.Vec3(0.5, -0.5, -0.5)
vertex3.normal = hg.Vec3(0, 0, -1)
vertex3.uv0 = hg.Vec2(1, 0)
d = mdl_builder:AddVertex(vertex3)

mdl_builder:AddTriangle(d, c, b) -- Adding the first triangle of the current face
mdl_builder:AddTriangle(b, a, d) -- Second triangle

-- +
vertex0 = hg.Vertex()
vertex0.pos = hg.Vec3(-0.5, -0.5, 0.5)
vertex0.normal = hg.Vec3(0, 0, 1)
vertex0.uv0 = hg.Vec2(0, 0)
a = mdl_builder:AddVertex(vertex0)

vertex1 = hg.Vertex()
vertex1.pos = hg.Vec3(-0.5,0.5, 0.5)
vertex1.normal = hg.Vec3(0, 0, 1)
vertex1.uv0 = hg.Vec2(0, 1)
b = mdl_builder:AddVertex(vertex1)

vertex2 = hg.Vertex()
vertex2.pos = hg.Vec3(0.5, 0.5, 0.5)
vertex2.normal = hg.Vec3(0, 0, 1)
vertex2.uv0 = hg.Vec2(1, 1)
c = mdl_builder:AddVertex(vertex2)

vertex3 = hg.Vertex()
vertex3.pos = hg.Vec3(0.5, -0.5, 0.5)
vertex3.normal = hg.Vec3(0, 0, 1)
vertex3.uv0 = hg.Vec2(1, 0)
d = mdl_builder:AddVertex(vertex3)

mdl_builder:AddTriangle(a, b, c)
mdl_builder:AddTriangle(a, c, d)

-- -
vertex0 = hg.Vertex()
vertex0.pos = hg.Vec3(-0.5, -0.5, -0.5)
vertex0.normal = hg.Vec3(0, -1, 0)
vertex0.uv0 = hg.Vec2(0, 0)
a = mdl_builder:AddVertex(vertex0)

vertex1 = hg.Vertex()
vertex1.pos = hg.Vec3(-0.5, -0.5, 0.5)
vertex1.normal = hg.Vec3(0, -1, 0)
vertex1.uv0 = hg.Vec2(0, 1)
b = mdl_builder:AddVertex(vertex1)

vertex2 = hg.Vertex()
vertex2.pos = hg.Vec3(0.5, -0.5, 0.5)
vertex2.normal = hg.Vec3(0, -1, 0)
vertex2.uv0 = hg.Vec2(1, 1)
c = mdl_builder:AddVertex(vertex2)

vertex3 = hg.Vertex()
vertex3.pos = hg.Vec3(0.5, -0.5, -0.5)
vertex3.normal = hg.Vec3(0, -1, 0)
vertex3.uv0 = hg.Vec2(1, 0)
d = mdl_builder:AddVertex(vertex3)

mdl_builder:AddTriangle(a, b, c)
mdl_builder:AddTriangle(a, c, d)

-- +
vertex0 = hg.Vertex()
vertex0.pos = hg.Vec3(-0.5, 0.5, -0.5)
vertex0.normal = hg.Vec3(0, 1, 0)
vertex0.uv0 = hg.Vec2(0, 0)
a = mdl_builder:AddVertex(vertex0)

vertex1 = hg.Vertex()
vertex1.pos = hg.Vec3(-0.5, 0.5, 0.5)
vertex1.normal = hg.Vec3(0, 1, 0)
vertex1.uv0 = hg.Vec2(0, 1)
b = mdl_builder:AddVertex(vertex1)

vertex2 = hg.Vertex()
vertex2.pos = hg.Vec3(0.5, 0.5, 0.5)
vertex2.normal = hg.Vec3(0, 1, 0)
vertex2.uv0 = hg.Vec2(1, 1)
c = mdl_builder:AddVertex(vertex2)

vertex3 = hg.Vertex()
vertex3.pos = hg.Vec3(0.5, 0.5, -0.5)
vertex3.normal = hg.Vec3(0, 1, 0)
vertex3.uv0 = hg.Vec2(1, 0)
d = mdl_builder:AddVertex(vertex3)

mdl_builder:AddTriangle(d, c, b)
mdl_builder:AddTriangle(b, a, d)

-- -
vertex0 = hg.Vertex()
vertex0.pos = hg.Vec3(-0.5, -0.5, -0.5)
vertex0.normal = hg.Vec3(-1, 0, 0)
vertex0.uv0 = hg.Vec2(0, 0)
a = mdl_builder:AddVertex(vertex0)

vertex1 = hg.Vertex()
vertex1.pos = hg.Vec3(-0.5, -0.5, 0.5)
vertex1.normal = hg.Vec3(-1, 0, 0)
vertex1.uv0 = hg.Vec2(0, 1)
b = mdl_builder:AddVertex(vertex1)

vertex2 = hg.Vertex()
vertex2.pos = hg.Vec3(-0.5, 0.5, 0.5)
vertex2.normal = hg.Vec3(-1, 0, 0)
vertex2.uv0 = hg.Vec2(1, 1)
c = mdl_builder:AddVertex(vertex2)

vertex3 = hg.Vertex()
vertex3.pos = hg.Vec3(-0.5, 0.5, -0.5)
vertex3.normal = hg.Vec3(-1, 0, 0)
vertex3.uv0 = hg.Vec2(1, 0)
d = mdl_builder:AddVertex(vertex3)

mdl_builder:AddTriangle(d, c, b)
mdl_builder:AddTriangle(b, a, d)

-- +
vertex0 = hg.Vertex()
vertex0.pos = hg.Vec3(0.5, -0.5, -0.5)
vertex0.normal = hg.Vec3(1, 0, 0)
vertex0.uv0 = hg.Vec2(0, 0)
a = mdl_builder:AddVertex(vertex0)

vertex1 = hg.Vertex()
vertex1.pos = hg.Vec3(0.5, -0.5, 0.5)
vertex1.normal = hg.Vec3(1, 0, 0)
vertex1.uv0 = hg.Vec2(0, 1)
b = mdl_builder:AddVertex(vertex1)

vertex2 = hg.Vertex()
vertex2.pos = hg.Vec3(0.5, 0.5, 0.5)
vertex2.normal = hg.Vec3(1, 0, 0)
vertex2.uv0 = hg.Vec2(1, 1)
c = mdl_builder:AddVertex(vertex2)

vertex3 = hg.Vertex()
vertex3.pos = hg.Vec3(0.5, 0.5, -0.5)
vertex3.normal = hg.Vec3(1, 0, 0)
vertex3.uv0 = hg.Vec2(1, 0)
d = mdl_builder:AddVertex(vertex3)

mdl_builder:AddTriangle(a, b, c)
mdl_builder:AddTriangle(a, c, d)

cube_mdl = mdl_builder:MakeModel(vtx_layout) -- Create the actuel cube model

ground_mdl = hg.CreatePlaneModel(vtx_layout, 5, 5, 1, 1)

shader = hg.LoadProgramFromFile('resources_compiled/shaders/mdl')

-- main loop
angle = 0

while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(win) do
    dt = hg.TickClock()
    angle = angle + hg.time_to_sec_f(dt)

    viewpoint = hg.TranslationMat4(hg.Vec3(0, 1, -3))
    hg.SetViewPerspective(0, 0, 0, res_x, res_y, viewpoint)

    hg.DrawModel(0, cube_mdl, shader, {}, {}, hg.TransformationMat4(
        hg.Vec3(0, 1, 0), hg.Vec3(angle, angle, angle)))
    hg.DrawModel(0, ground_mdl, shader, {}, {},
                 hg.TranslationMat4(hg.Vec3(0, 0, 0)))

    hg.Frame()
    hg.UpdateWindow(win)
end

hg.RenderShutdown()
