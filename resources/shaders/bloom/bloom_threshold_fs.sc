$input v_texcoord0

#include <bgfx_shader.sh>

uniform vec4 u_params;
SAMPLER2D(u_source, 0);

#define u_threshold u_params.x
#define u_knee u_params.y

void main() {
	vec4 color = texture2D(u_source, v_texcoord0);
	float lum = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
	float r = clamp(lum - u_threshold + u_knee, 0, 2.0 * u_knee);
	r = (r * r) / (4.0 * u_knee);
	gl_FragColor = color * max(r , lum - u_threshold) / max(lum, 0.00001);
}
