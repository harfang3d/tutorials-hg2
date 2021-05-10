$input v_viewRay

#include <forward_pipeline.sh>

SAMPLER2D(u_color, 0);
SAMPLER2D(u_attr0, 1);
SAMPLER2D(u_attr1, 2);
SAMPLER2D(u_noise, 3);
SAMPLERCUBE(u_probe, 4);
SAMPLER2D(u_depthTex, 5); // input: minimum depth pyramid

#define sample_count 2

#define uv_ratio vec2_splat(uAAAParams[0].x)

#include <aaa_utils.sh>
#include <hiz_trace.sh>

void main() {
	vec4 jitter = texture2D(u_noise, mod(gl_FragCoord.xy, vec2(64, 64)) / vec2(64, 64));

	// sample normal/depth
	vec2 uv = gl_FragCoord.xy / uResolution.xy;
	vec4 attr0 = texture2D(u_attr0, uv * uv_ratio);

	vec3 n = attr0.xyz;

	// compute ray origin & direction
	vec3 ray_o = v_viewRay * attr0.w;

	// spread
	vec3 right = normalize(cross(n, vec3(1, 0, 0)));
	vec3 up = cross(n, right);

	//
	vec4 color = vec4_splat(0.);

	for (int i = 0; i < sample_count; ++i) {
		float r = (float(i) + jitter.y) / float(sample_count);
		float spread = r * 3.141592 * 0.5 * 0.9;
		float cos_spread = cos(spread), sin_spread = sin(spread);

		for (int j = 0; j < sample_count; ++j) {
			float angle = float(j + jitter.w) / float(sample_count) * 3.141592 * 2.0;
			float cos_angle = cos(angle), sin_angle = sin(angle);
			vec3 ray_d_spread = (right * cos_angle + up * sin_angle) * sin_spread + n * cos_spread;

			vec2 hit_pixel;
			vec3 hit_point;
			if (TraceScreenRay(ray_o - v_viewRay * 0.0, ray_d_spread, uMainProjection, 0.1, /*jitter.z,*/ 64, hit_pixel, hit_point)) {
				// use hit pixel velocity to compensate the fact that we are sampling the previous frame
				vec2 uv = hit_pixel / uResolution.xy;
				vec2 vel = GetVelocityVector(uv);

				vec4 attr0 = texelFetch(u_attr0, ivec2(hit_pixel), 0);

				if (dot(attr0.xyz, ray_d_spread) < 0.0) { // ray facing the collision
					vec3 irradiance = texture2D(u_color, uv - vel * uv_ratio).xyz;
					color += vec4(irradiance * 1.0, 0.0);
				} else {
					color += vec4(0.0, 0.0, 0.0, 0.0); // backface hit
				}
			} else {
				vec3 world_ray_d_spread = mul(uMainInvView, vec4(ray_d_spread, 0.0)).xyz;
				color += vec4(textureCubeLod(u_probe, world_ray_d_spread, 0).xyz, 1.);
			}
		}
	}

#if 1
	if (isNan(color.x) || isNan(color.y) || isNan(color.z))
		color = vec4(1., 0., 0., 1.);
#endif

	gl_FragColor = color / float(sample_count * sample_count);
}
