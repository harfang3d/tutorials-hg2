$input a_position, a_normal, a_texcoord0
$output vNormal


#include <bgfx_shader.sh>

void main() {
	vNormal = mul(u_model[0], vec4(a_normal * 2.0 - 1.0, 0.0)).xyz;
	gl_Position = mul(u_modelViewProj, vec4(a_position, 1.0));
}
