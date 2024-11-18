-- Toyota 2JZ-GTE Engine model by Serhii Denysenko (CGTrader: serhiidenysenko8256)
-- URL : https://www.cgtrader.com/3d-models/vehicle/part/toyota-2jz-gte-engine-2932b715-2f42-4ecd-93ce-df9507c67ce8

hg = require("harfang")
hg.InputInit()
hg.WindowSystemInit()

local res_x, res_y, tex_size = 1024, 1024, 1024
win = hg.RenderInit('Frame Buffer Screenshot', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

--Link precompiled assets folder to the project
hg.AddAssetsFolder("resources_compiled")

--Create Pipeline
pipeline = hg.CreateForwardPipeline()
local res = hg.PipelineResources()

--Create a 1024*1024 frame buffer to draw the scene to
local frame_buffer = hg.CreateFrameBuffer(tex_size, tex_size, hg.TF_RGBA8, hg.TF_D24, 4, 'framebuffer')
local tex_color = hg.GetColorTexture(frame_buffer)

--Prepare screenshot requirements
local tex_color_ref = res:AddTexture("tex_rb", tex_color)
local tex_readback = hg.CreateTexture(tex_size, tex_size, "readback", hg.TF_ReadBack | hg.TF_BlitDestination, hg.TF_RGBA8)
local picture = hg.Picture(tex_size, tex_size, hg.PF_RGBA32)

--Load scene
scene = hg.Scene()
ret = hg.LoadSceneFromAssets("car_engine/engine.scn", scene, res, hg.GetForwardPipelineInfo())

--Create the plane model
vtx_layout = hg.VertexLayoutPosFloatTexCoord0UInt8()
plane_mdl = hg.CreatePlaneModel(vtx_layout, 1, 1, 1, 1)
plane_ref = res:AddModel('plane', plane_mdl)

--Prepare the plane shader program
plane_prg = hg.LoadProgramFromAssets('shaders/texture')

--Main Loop 
local frame = 0
local view_id = 0
local state = "none"

while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(win) do
	dt = hg.TickClock()

    --Change state to capture when user press Space and Capture is not already running
    if(hg.ReadKeyboard():Key(hg.K_Space) and state == "none") then
        state = "capture"
        frame_count_capture, view_id = hg.CaptureTexture(view_id, res, tex_color_ref, tex_readback, picture)
        print(frame_count_capture, frame) 

    --Take screenshot if CaptureTexture is ready and user pressed space
    elseif(state == "capture" and frame_count_capture <= frame) then
        local png_filename = "screenshots/" .. ".png"
        hg.SavePNG(picture, png_filename)
        state = "none" --Reset state to none to be able to screenshot again
    end

    --Update Scene and render to the frameBuffer
	scene:Update(dt)

    view_id = 0
	view_id = hg.SubmitSceneToPipeline(view_id, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res, frame_buffer.handle)
	
    --Draw a plabe using the texture the scene was rendered to
	hg.SetViewPerspective(view_id, 0, 0, res_x, res_y, hg.TranslationMat4(hg.Vec3(0, 0, -1.8)))
	
	val_uniforms = {hg.MakeUniformSetValue('color', hg.Vec4(1, 1, 1, 1))} 
	tex_uniforms = {hg.MakeUniformSetTexture('s_tex', tex_color, 0)}
	
	hg.DrawModel(view_id, plane_mdl, plane_prg, val_uniforms, tex_uniforms, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Vec3(math.pi / 2, 0, math.pi)))
    
    --
	frame = hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)


  
