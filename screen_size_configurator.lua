hg = require("harfang")

-- Init input and window system
hg.InputInit()
hg.WindowSystemInit()

res_list = {{640, 360}, {768, 432}, {896, 504}, {1024, 576}, {1152, 648}, {1280, 720}, {1920, 1080}, {1920, 1200}, {2560, 1440}, {3840, 2160}, {5120, 2880}}
res_list_str = {}
-- Prepare list of resolutions
local i
for i = 1, #res_list do
    table.insert(res_list_str, res_list[i][1] .. "x" .. res_list[i][2])
end

-- Window mode selection
mode_list = {hg.WV_Windowed, hg.WV_Fullscreen, hg.WV_Undecorated, hg.WV_FullscreenMonitor1, hg.WV_FullscreenMonitor2, hg.WV_FullscreenMonitor3}
mode_list_str = {"Windowed", "Fullscreen", "Undecorated", "Fullscreen Monitor #1", "Fullscreen Monitor #2", "Fullscreen Monitor #3"}
-- Resize the lists based on the available monitors
mon_list = hg.GetMonitors()
mon_list_size = mon_list:size()
-- If only 2 monitors, remove the third
if mon_list_size == 2 then
    mode_list = {table.unpack(mode_list, 1, 5)}
    mode_list_str = {table.unpack(mode_list_str, 1, 5)}
end
-- If only 1 monitor, remove the second and the third
if mon_list_size == 1 then
    mode_list = {table.unpack(mode_list, 1, 4)}
    mode_list_str = {table.unpack(mode_list_str, 1, 4)}
end

-- Preset values define by the configurator
res_preset = 5
window_mode_preset = 0

-- -- Configurator window --
-- Init configurator window
config_res_x, config_res_y = 600, 200
config_window_mode = hg.WV_Windowed

config_win = hg.NewWindow("Window Configurator", config_res_x, config_res_y, 32, config_window_mode)
hg.RenderInit(config_win)
hg.RenderReset(config_res_x, config_res_y, hg.RF_VSync)

-- Init shader
hg.AddAssetsFolder('resources_compiled')
imgui_prg = hg.LoadProgramFromAssets('core/shader/imgui')
imgui_img_prg = hg.LoadProgramFromAssets('core/shader/imgui_image')

-- Init imgui
hg.ImGuiInit(10, imgui_prg, imgui_img_prg)

-- Configurator window main loop
press_apply = False
while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(config_win) and not press_apply do
    hg.ImGuiBeginFrame(config_res_x, config_res_y, hg.TickClock(), hg.ReadMouse(), hg.ReadKeyboard())

    -- Imgui window configutator 
    if hg.ImGuiBegin("Window Configuration", true, hg.ImGuiWindowFlags_NoMove | hg.ImGuiWindowFlags_NoResize) then
        -- Set the ui position and size in the window
        hg.ImGuiSetWindowPos("Window Configuration", hg.Vec2(0, 0), hg.ImGuiCond_Once)
        hg.ImGuiSetWindowSize("Window Configuration", hg.Vec2(config_res_x, config_res_y), hg.ImGuiCond_Once)
        
        -- Create the 2 combo list to chose the resolution and the window mode
        hg.ImGuiText("Screen")
        res_modified, res_preset = hg.ImGuiCombo("Resolution", res_preset, res_list_str)
        fullscreen_modified, window_mode_preset = hg.ImGuiCombo("Mode", window_mode_preset, mode_list_str)

        -- Apply button
        hg.ImGuiSpacing()
        hg.ImGuiPushStyleColor(hg.ImGuiCol_Button, hg.Color(0.0, 0.5, 1.0, 1.0))
        press_apply = hg.ImGuiButton("Apply")
        hg.ImGuiPopStyleColor()
    end

    hg.ImGuiEnd()

    hg.SetView2D(0, 0, 0, config_res_x, config_res_y, -1, 1, hg.CF_Color | hg.CF_Depth, hg.Color.Black, 1, 0)
    hg.ImGuiEndFrame(0)

    hg.Frame()
    hg.UpdateWindow(config_win)
end

hg.RenderShutdown()
hg.DestroyWindow(config_win)

-- -- Preset window --
if press_apply then
    -- Init preset window
    selected_res = res_list[res_preset + 1]
    selected_res_x, selected_res_y = selected_res[1], selected_res[2]
    selected_window_mode = mode_list[window_mode_preset + 1]

    win = hg.NewWindow("Window preset", selected_res_x, selected_res_y, 32, selected_window_mode)
    hg.RenderInit(win)
    hg.RenderReset(selected_res_x, selected_res_y, hg.RF_VSync)

    -- Init shader
    imgui_prg = hg.LoadProgramFromAssets('core/shader/imgui')
    imgui_img_prg = hg.LoadProgramFromAssets('core/shader/imgui_image')

    -- Init imgui
    hg.ImGuiInit(10, imgui_prg, imgui_img_prg)
    imgui_window_size = hg.Vec2(190, 50)
    imgui_window_pos = hg.Vec2((selected_res_x/2)-(imgui_window_size.x/2), (selected_res_y/2)-(imgui_window_size.y/2))

    -- Preset window main loop
    while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(win) do
        hg.ImGuiBeginFrame(selected_res_x, selected_res_y, hg.TickClock(), hg.ReadMouse(), hg.ReadKeyboard())

        -- Imgui window configutator 
        if hg.ImGuiBegin("Window Parameters", true, hg.ImGuiWindowFlags_NoResize | hg.ImGuiWindowFlags_NoTitleBar) then
            -- Set the ui position and size in the window
            hg.ImGuiSetWindowPos("Window Parameters", imgui_window_pos, hg.ImGuiCond_Once)
            hg.ImGuiSetWindowSize("Window Parameters", imgui_window_size, hg.ImGuiCond_Once)

            -- Display window parameters selected in the configurator
            hg.ImGuiText(string.format("Window size : %d x %d", selected_res_x, selected_res_y))
            hg.ImGuiText(string.format("Window mode : %s", mode_list_str[window_mode_preset + 1]))
        end

        hg.ImGuiEnd()

        hg.SetView2D(0, 0, 0, selected_res_x, selected_res_y, -1, 1, hg.CF_Color | hg.CF_Depth, hg.Color.Red/2, 1, 0)
        hg.ImGuiEndFrame(0)

        hg.Frame()
        hg.UpdateWindow(win)
    end

    hg.RenderShutdown()
    hg.DestroyWindow(win)
end