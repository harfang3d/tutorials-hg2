$input v_texcoord0

#include <bgfx_shader.sh>

SAMPLER2D(s_tex, 0);

void main() {
	float range_low = 0.998, range_high = 0.9985;
	float depth = (texture2D(s_tex, v_texcoord0).r - range_low) / (range_high - range_low);
	gl_FragColor = vec4(depth, depth, depth, 1.0);
}
