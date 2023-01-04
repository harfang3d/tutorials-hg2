// Create and draw models without a pipeline and using ModelBuilder

package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
    hg.InputInit()
    hg.WindowSystemInit()

    res_x, res_y := int32(1280), int32(720)
    win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Draw and Create Models using ModelBuilder - no Pipeline",
                        res_x, res_y, hg.RFVSync | hg.RFMSAA4X)

    // vertex layout, materials and models
    vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

    mdl_builder := hg.NewModelBuilder()

    // Below are the 6 faces of the cube model, with 4 vertices declared for each face, that are added as 2 triangles to make a complete quad face

    // -
    vertex0 := hg.NewVertex() // Vertex constructor
    vertex0.SetPos(hg.NewVec3WithXYZ(-0.5, -0.5, -0.5)) // Declaring the vertex position
    vertex0.SetNormal(hg.NewVec3WithXYZ(0, 0, -1)) // Declaring the normal
    vertex0.SetUv0(hg.NewVec2WithXY(0, 0)) // Declaring the UV
    a := mdl_builder.AddVertex(vertex0) // Adding the vertex to the ModelBuilder

    vertex1 := hg.NewVertex()
    vertex1.SetPos(hg.NewVec3WithXYZ(-0.5, 0.5, -0.5))
    vertex1.SetNormal(hg.NewVec3WithXYZ(0, 0, -1))
    vertex1.SetUv0(hg.NewVec2WithXY(0, 1))
    b := mdl_builder.AddVertex(vertex1)

    vertex2 := hg.NewVertex()
    vertex2.SetPos(hg.NewVec3WithXYZ(0.5, 0.5, -0.5))
    vertex2.SetNormal(hg.NewVec3WithXYZ(0, 0, -1))
    vertex2.SetUv0(hg.NewVec2WithXY(1, 1))
    c := mdl_builder.AddVertex(vertex2)

    vertex3 := hg.NewVertex()
    vertex3.SetPos(hg.NewVec3WithXYZ(0.5, -0.5, -0.5))
    vertex3.SetNormal(hg.NewVec3WithXYZ(0, 0, -1))
    vertex3.SetUv0(hg.NewVec2WithXY(1, 0))
    d := mdl_builder.AddVertex(vertex3)

    mdl_builder.AddTriangle(d, c, b) // Adding the first triangle of the current face
    mdl_builder.AddTriangle(b, a, d) // Second triangle

    // +
    vertex0 = hg.NewVertex()
    vertex0.SetPos(hg.NewVec3WithXYZ(-0.5, -0.5, 0.5))
    vertex0.SetNormal(hg.NewVec3WithXYZ(0, 0, 1))
    vertex0.SetUv0(hg.NewVec2WithXY(0, 0))
    a = mdl_builder.AddVertex(vertex0)

    vertex1 = hg.NewVertex()
    vertex1.SetPos(hg.NewVec3WithXYZ(-0.5,0.5, 0.5))
    vertex1.SetNormal(hg.NewVec3WithXYZ(0, 0, 1))
    vertex1.SetUv0(hg.NewVec2WithXY(0, 1))
    b = mdl_builder.AddVertex(vertex1)

    vertex2 = hg.NewVertex()
    vertex2.SetPos(hg.NewVec3WithXYZ(0.5, 0.5, 0.5))
    vertex2.SetNormal(hg.NewVec3WithXYZ(0, 0, 1))
    vertex2.SetUv0(hg.NewVec2WithXY(1, 1))
    c = mdl_builder.AddVertex(vertex2)

    vertex3 = hg.NewVertex()
    vertex3.SetPos(hg.NewVec3WithXYZ(0.5, -0.5, 0.5))
    vertex3.SetNormal(hg.NewVec3WithXYZ(0, 0, 1))
    vertex3.SetUv0(hg.NewVec2WithXY(1, 0))
    d = mdl_builder.AddVertex(vertex3)

    mdl_builder.AddTriangle(a, b, c)
    mdl_builder.AddTriangle(a, c, d)

    // -
    vertex0 = hg.NewVertex()
    vertex0.SetPos(hg.NewVec3WithXYZ(-0.5, -0.5, -0.5))
    vertex0.SetNormal(hg.NewVec3WithXYZ(0, -1, 0))
    vertex0.SetUv0(hg.NewVec2WithXY(0, 0))
    a = mdl_builder.AddVertex(vertex0)

    vertex1 = hg.NewVertex()
    vertex1.SetPos(hg.NewVec3WithXYZ(-0.5, -0.5, 0.5))
    vertex1.SetNormal(hg.NewVec3WithXYZ(0, -1, 0))
    vertex1.SetUv0(hg.NewVec2WithXY(0, 1))
    b = mdl_builder.AddVertex(vertex1)

    vertex2 = hg.NewVertex()
    vertex2.SetPos(hg.NewVec3WithXYZ(0.5, -0.5, 0.5))
    vertex2.SetNormal(hg.NewVec3WithXYZ(0, -1, 0))
    vertex2.SetUv0(hg.NewVec2WithXY(1, 1))
    c = mdl_builder.AddVertex(vertex2)

    vertex3 = hg.NewVertex()
    vertex3.SetPos(hg.NewVec3WithXYZ(0.5, -0.5, -0.5))
    vertex3.SetNormal(hg.NewVec3WithXYZ(0, -1, 0))
    vertex3.SetUv0(hg.NewVec2WithXY(1, 0))
    d = mdl_builder.AddVertex(vertex3)

    mdl_builder.AddTriangle(a, b, c)
    mdl_builder.AddTriangle(a, c, d)

    // +
    vertex0 = hg.NewVertex()
    vertex0.SetPos(hg.NewVec3WithXYZ(-0.5, 0.5, -0.5))
    vertex0.SetNormal(hg.NewVec3WithXYZ(0, 1, 0))
    vertex0.SetUv0(hg.NewVec2WithXY(0, 0))
    a = mdl_builder.AddVertex(vertex0)

    vertex1 = hg.NewVertex()
    vertex1.SetPos(hg.NewVec3WithXYZ(-0.5, 0.5, 0.5))
    vertex1.SetNormal(hg.NewVec3WithXYZ(0, 1, 0))
    vertex1.SetUv0(hg.NewVec2WithXY(0, 1))
    b = mdl_builder.AddVertex(vertex1)

    vertex2 = hg.NewVertex()
    vertex2.SetPos(hg.NewVec3WithXYZ(0.5, 0.5, 0.5))
    vertex2.SetNormal(hg.NewVec3WithXYZ(0, 1, 0))
    vertex2.SetUv0(hg.NewVec2WithXY(1, 1))
    c = mdl_builder.AddVertex(vertex2)

    vertex3 = hg.NewVertex()
    vertex3.SetPos(hg.NewVec3WithXYZ(0.5, 0.5, -0.5))
    vertex3.SetNormal(hg.NewVec3WithXYZ(0, 1, 0))
    vertex3.SetUv0(hg.NewVec2WithXY(1, 0))
    d = mdl_builder.AddVertex(vertex3)

    mdl_builder.AddTriangle(d, c, b)
    mdl_builder.AddTriangle(b, a, d)

    // -
    vertex0 = hg.NewVertex()
    vertex0.SetPos(hg.NewVec3WithXYZ(-0.5, -0.5, -0.5))
    vertex0.SetNormal(hg.NewVec3WithXYZ(-1, 0, 0))
    vertex0.SetUv0(hg.NewVec2WithXY(0, 0))
    a = mdl_builder.AddVertex(vertex0)

    vertex1 = hg.NewVertex()
    vertex1.SetPos(hg.NewVec3WithXYZ(-0.5, -0.5, 0.5))
    vertex1.SetNormal(hg.NewVec3WithXYZ(-1, 0, 0))
    vertex1.SetUv0(hg.NewVec2WithXY(0, 1))
    b = mdl_builder.AddVertex(vertex1)

    vertex2 = hg.NewVertex()
    vertex2.SetPos(hg.NewVec3WithXYZ(-0.5, 0.5, 0.5))
    vertex2.SetNormal(hg.NewVec3WithXYZ(-1, 0, 0))
    vertex2.SetUv0(hg.NewVec2WithXY(1, 1))
    c = mdl_builder.AddVertex(vertex2)

    vertex3 = hg.NewVertex()
    vertex3.SetPos(hg.NewVec3WithXYZ(-0.5, 0.5, -0.5))
    vertex3.SetNormal(hg.NewVec3WithXYZ(-1, 0, 0))
    vertex3.SetUv0(hg.NewVec2WithXY(1, 0))
    d = mdl_builder.AddVertex(vertex3)

    mdl_builder.AddTriangle(d, c, b)
    mdl_builder.AddTriangle(b, a, d)

    // +
    vertex0 = hg.NewVertex()
    vertex0.SetPos(hg.NewVec3WithXYZ(0.5, -0.5, -0.5))
    vertex0.SetNormal(hg.NewVec3WithXYZ(1, 0, 0))
    vertex0.SetUv0(hg.NewVec2WithXY(0, 0))
    a = mdl_builder.AddVertex(vertex0)

    vertex1 = hg.NewVertex()
    vertex1.SetPos(hg.NewVec3WithXYZ(0.5, -0.5, 0.5))
    vertex1.SetNormal(hg.NewVec3WithXYZ(1, 0, 0))
    vertex1.SetUv0(hg.NewVec2WithXY(0, 1))
    b = mdl_builder.AddVertex(vertex1)

    vertex2 = hg.NewVertex()
    vertex2.SetPos(hg.NewVec3WithXYZ(0.5, 0.5, 0.5))
    vertex2.SetNormal(hg.NewVec3WithXYZ(1, 0, 0))
    vertex2.SetUv0(hg.NewVec2WithXY(1, 1))
    c = mdl_builder.AddVertex(vertex2)

    vertex3 = hg.NewVertex()
    vertex3.SetPos(hg.NewVec3WithXYZ(0.5, 0.5, -0.5))
    vertex3.SetNormal(hg.NewVec3WithXYZ(1, 0, 0))
    vertex3.SetUv0(hg.NewVec2WithXY(1, 0))
    d = mdl_builder.AddVertex(vertex3)

    mdl_builder.AddTriangle(a, b, c)
    mdl_builder.AddTriangle(a, c, d)

    cube_mdl := mdl_builder.MakeModel(vtx_layout) // Create the actual cube model

    ground_mdl := hg.CreatePlaneModel(vtx_layout, 5, 5, 1, 1)

    shader := hg.LoadProgramFromFile("resources_compiled/shaders/mdl")

    // main loop
    angle := float32(0.0)

    for !hg.ReadKeyboardWithName("raw").Key(hg.KEscape) {
        dt := hg.TickClock()
		angle = angle + hg.TimeToSecF(dt)

        viewpoint := hg.TranslationMat4(hg.NewVec3WithXYZ(0, 1, -3))
        hg.SetViewPerspective(0, 0, 0, res_x, res_y, viewpoint)

		hg.DrawModelWithSliceOfValuesSliceOfTextures(0, cube_mdl, shader, hg.GoSliceOfUniformSetValue{}, hg.GoSliceOfUniformSetTexture{}, hg.TransformationMat4(hg.NewVec3WithXYZ(0.0, 1, 0.0), hg.NewVec3WithXYZ(angle, angle, angle)))
		hg.DrawModelWithSliceOfValuesSliceOfTextures(0, ground_mdl, shader, hg.GoSliceOfUniformSetValue{}, hg.GoSliceOfUniformSetTexture{}, hg.TranslationMat4(hg.NewVec3WithXYZ(0.0, 0.0, 0.0)))


        hg.Frame()
        hg.UpdateWindow(win)
    }
    hg.RenderShutdown()
}