$input v_texcoord0

#include <bgfx_shader.sh>

SAMPLER2D(u_source, 0);
SAMPLER2D(u_input, 1);

void main() {
	vec3 color = texture2D(u_source, v_texcoord0.xy).rgb;
	color += texture2D(u_input, v_texcoord0.xy).rgb;
	
	// [todo] tone mapping ?
	// [todo] gamme correction ?
	// [todo] color = clamp(color, 0.0, 1.0);

	gl_FragColor = vec4(color, 1.0);
}
