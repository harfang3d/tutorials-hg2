$input v_texcoord0

#include <bgfx_shader.sh>

SAMPLER2D(u_source, 0);
uniform vec4 u_params;

void main() {
	vec4 far_offset = vec4(-u_params.x, u_params.x, u_params.y, 0.0);
	vec4 offset = 0.5*far_offset;
	vec2 uv = v_texcoord0.xy;
		
	vec4 s0 = texture2D(u_source, uv - far_offset.yz);
	vec4 s1 = texture2D(u_source, uv - far_offset.wz);
	vec4 s2 = texture2D(u_source, uv - far_offset.xz);

	vec4 s3 = texture2D(u_source, uv + far_offset.xw);
	vec4 s4 = texture2D(u_source, uv);
	vec4 s5 = texture2D(u_source, uv + far_offset.yw);

	vec4 s6 = texture2D(u_source, uv + far_offset.xz);
	vec4 s7 = texture2D(u_source, uv + far_offset.wz);
	vec4 s8 = texture2D(u_source, uv + far_offset.yz);
	
	vec4 t0 = texture2D(u_source, uv - offset.yz);
	vec4 t1 = texture2D(u_source, uv - offset.xz);
	vec4 t2 = texture2D(u_source, uv + offset.xz);
	vec4 t3 = texture2D(u_source, uv + offset.yz);
	
	
	vec4 v0 = s0 + s1 + s3 + s4;
	vec4 v1 = s1 + s2 + s4 + s5;
	vec4 v2 = s3 + s4 + s6 + s7;
	vec4 v3 = s4 + s5 + s7 + s8;
	vec4 v4 = t0 + t1 + t2 + t3;

	gl_FragColor = (((v0 + v1 + v2 + v3) / 4.0) + v4) / 8.0;
}
