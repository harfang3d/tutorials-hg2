$input vNormal

#include <bgfx_shader.sh>

void main() {
	vec3 light = vec3(1,0.8,0.5);
	
	vec4 color = vec4(1.,1.,1.,1.);
	vec4 ambient_color = vec4(0.1,0.1,0.2,1.);
	
	vec4 backlight_color=vec4(0.75,0.85,1.,1.);
	vec4 mainlight_color=vec4(1.,1.,1.,1.);
	
	float backlight_intensity = 0.5;
	light = normalize(light);
	vec3 normal = normalize(vNormal);
	
	float main_lighting = max(0.,-dot(normal,light));
	float back_lighting = max(0.,-dot(normal,-light)) * backlight_intensity;
	
	
	vec4 light_color = min(color*(mainlight_color*main_lighting + backlight_color * back_lighting) + ambient_color, vec4(1.,1.,1.,1.));
	
	gl_FragColor = light_color;
}
