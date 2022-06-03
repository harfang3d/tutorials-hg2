-- Instantiating scenes

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Scene instances', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder("resources_compiled")

-- rendering pipeline
pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

-- load host scene
scene = hg.Scene()
hg.LoadSceneFromAssets('playground/playground.scn', scene, res, hg.GetForwardPipelineInfo())


-- declare the biped actor class
BipedActor =
{
	__node = nil,
	__delay = 0,
	__state = nil,
	__playing_anim_ref = nil
}

function BipedActor:new(pos)
	local o={}
	setmetatable(o,self)
	self.__index=self
	o.__node = hg.CreateInstanceFromAssets(scene, hg.Mat4.Identity, "biped/biped.scn", res, hg.GetForwardPipelineInfo())
	o.__node:GetTransform():SetPosRot(pos, hg.Deg3(0, hg.FRand(360), 0))
	o.__delay = 0
	o.__state = nil
	o.__playing_anim_ref = nil
	return o
end

function BipedActor:__start_anim(name)
	anim = self.__node:GetInstanceSceneAnim(name)  -- get instance specific animation
	if self.__playing_anim_ref ~= nil then
		scene:StopAnim(self.__playing_anim_ref)
	end
	self.__playing_anim_ref = scene:PlayAnim(anim, hg.ALM_Loop)
end

function BipedActor:update(dt)
	-- check for state change
	self.__delay = self.__delay - dt

	if self.__delay <= 0 then
		states = {'idle', 'walk', 'run'}
		self.__state = states[hg.Rand(#states) + 1]
		self.__delay = self.__delay + hg.time_from_sec_f(hg.FRRand(2, 6))  -- 2 to 6 seconds before next state change
		self:__start_anim(self.__state)
	end

	-- apply motion
	dt_sec_f = hg.time_to_sec_f(dt)

	transform = self.__node:GetTransform()
	pos, rot = transform:GetPosRot()

	if self.__state == 'walk' then
		pos = pos - hg.GetZ(transform:GetWorld()) * hg.Mtr(1.15) * dt_sec_f  -- 1.15 m/sec
		rot.y = rot.y + hg.Deg(50) * dt_sec_f
	elseif self.__state == 'run' then
		pos = pos - hg.GetZ(transform:GetWorld()) * hg.Mtr(4.5) * dt_sec_f  -- 4.5 m/sec
		rot.y = rot.y - hg.Deg(70) * dt_sec_f
	end

	-- confine actor to playground
	pos = hg.Clamp(pos, hg.Vec3(-10, 0, -10), hg.Vec3(10, 0, 10))

	transform:SetPosRot(pos, rot)
end

function BipedActor:destroy()
	scene:DestroyNode(self.__node)
end


-- spawn initial actors
actors = {}
for i=0, 19 do
	actor = BipedActor:new(hg.RandomVec3(hg.Vec3(-10, 0, -10), hg.Vec3(10, 0, 10)))
	table.insert(actors, actor)
end
print(string.format('%d nodes in scene', scene:GetAllNodeCount()))

-- main loop
keyboard = hg.Keyboard()

while not keyboard:Down(hg.K_Escape) and hg.IsWindowOpen(win) do
	_, res_x, res_y = hg.RenderResetToWindow(win, res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X | hg.RF_MaxAnisotropy)

	keyboard:Update()

	if keyboard:Pressed(hg.K_S) then
		table.insert(actors, BipedActor:new(hg.RandomVec3(hg.Vec3(-10, 0, -10), hg.Vec3(10, 0, 10))))
	end

	if keyboard:Pressed(hg.K_D) then
		if #actors > 0 then
			table.remove(actors,1):destroy()
			scene:GarbageCollect()
		end
	end

	dt = hg.TickClock()

	for i, actor in ipairs(actors) do
		actor:update(dt)
	end

	scene:Update(dt)

	view_state = hg.ComputePerspectiveViewState(hg.Mat4LookAt(hg.Vec3(0, 10, -14), hg.Vec3(0, 1, -4)), hg.Deg(45), 0.01, 1000, hg.ComputeAspectRatioX(res_x, res_y))
	vid, passId = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), view_state, pipeline, res)

	hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)
