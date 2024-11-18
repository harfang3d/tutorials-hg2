-- Toyota 2JZ-GTE Engine model by Serhii Denysenko (CGTrader: serhiidenysenko8256)
-- URL : https://www.cgtrader.com/3d-models/vehicle/part/toyota-2jz-gte-engine-2932b715-2f42-4ecd-93ce-df9507c67ce8

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1024, 1024
win = hg.RenderInit('Screenshot Buffer', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder("resources_compiled")

pipeline = hg.CreateForwardPipeline()
local res = hg.PipelineResources()

scene = hg.Scene()
hg.LoadSceneFromAssets("car_engine/engine.scn", scene, res, hg.GetForwardPipelineInfo())

local tex_size = 1024
local picture = hg.Picture(tex_size, tex_size, hg.PF_RGBA32F)
local frame_buffer = hg.CreateFrameBuffer(tex_size, tex_size, hg.TF_RGBA8, hg.TF_D24, 4, 'framebuffer')
local tex_color = hg.GetColorTexture(frame_buffer)
local tex_color_ref = res:AddTexture("tex_rb", tex_color)
local tex_readback = hg.CreateTexture(tex_size, tex_size, "readback", hg.TF_ReadBack | hg.TF_BlitDestination, hg.TF_RGBA8)

local state = "none"

-- main loop
local frame = 0
local view_id = 0

vtx_layout = hg.VertexLayoutPosFloatTexCoord0UInt8()

plane_mdl = hg.CreatePlaneModel(vtx_layout, 1, 1, 1, 1)
plane_ref = res:AddModel('plane', plane_mdl)

plane_prg = hg.LoadProgramFromAssets('shaders/texture')


while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(win) do
	dt = hg.TickClock()

	trs = scene:GetNode('engine_master'):GetTransform()
	trs:SetRot(trs:GetRot() + hg.Vec3(0, hg.Deg(15) * hg.time_to_sec_f(dt), 0))

    if(hg.ReadKeyboard():Key(hg.K_Space) and state == "none") then
        state = "capture"
        frame_count_capture, view_id = hg.CaptureTexture(view_id, res, tex_color_ref, tex_readback, picture)
        print(frame_count_capture, frame) 

    elseif(state == "capture" and frame_count_capture <= frame) then
        local png_filename = "screenshots/" .. ".png"
        hg.SavePNG(picture, png_filename)
        state = "none"
    end

	scene:Update(dt)

	view_id = 0
	view_id = hg.SubmitSceneToPipeline(view_id, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res, frame_buffer.handle)
	--view_id = hg.SubmitSceneToPipeline(view_id, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

	-- hg.SetViewPerspective(view_id, 0, 0, res_x, res_y, hg.TranslationMat4(hg.Vec3(0, 0, -1.8)))

	-- val_uniforms = {hg.MakeUniformSetValue('color', hg.Vec4(1, 1, 1, 1))} 
	-- tex_uniforms = {hg.MakeUniformSetTexture('s_tex', tex_color, 0)}

	-- hg.DrawModel(view_id, plane_mdl, plane_prg, val_uniforms, tex_uniforms, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Vec3(-math.pi/2,0, 0)))
	frame = hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)


  
