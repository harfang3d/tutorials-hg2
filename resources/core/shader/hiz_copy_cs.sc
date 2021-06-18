#include "bgfx_compute.sh"

// Initializes the level 0 of the minimum depth pyramid.

SAMPLER2D(u_depth, 0);
IMAGE2D_WR(u_depthTexOut, rg32f, 1);    // output: level 0 of the min/max depth pyramid

uniform vec4 u_zBounds;                 // near(x) far(y) unused(z,w)

NUM_THREADS(16, 16, 1)
void main() {
	ivec4 viewport = ivec4(u_viewRect);
	ivec2 coord = ivec2(gl_GlobalInvocationID.xy) + ivec2(u_viewRect.xy);

#if BGFX_SHADER_LANGUAGE_GLSL
		ivec2 tex_coord = ivec2(coord.x, textureSize(u_depth, 0).y - 1 - coord.y);
#else
		ivec2 tex_coord = coord;
#endif

	vec2 z;
	if(all(bvec4(greaterThanEqual(coord, viewport.xy), lessThan(coord, viewport.xy + viewport.zw)))) {
		z = texelFetch(u_depth, tex_coord, 0).ww;

		// [todo] use the backface depth buffer instead
		z.y += 0.1;

		// Store logarithmic depth
		z.x = u_zBounds.y * z.x / ((u_zBounds.y - u_zBounds.x) * z.x + u_zBounds.x*u_zBounds.y);
		z.y = u_zBounds.y * z.y / ((u_zBounds.y - u_zBounds.x) * z.y + u_zBounds.x*u_zBounds.y);
	} else {
		// Set pixels outside the viewport to the far clipping distance.
		z = vec2(1.0, 0.0);
	}

	imageStore(u_depthTexOut, coord, vec4(z.x, z.y, 0, 1));
}