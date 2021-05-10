#include <bgfx_shader.sh>

#define PI 3.14159265359

uniform vec4 uClock; // clock

// Environment
uniform vec4 uFogColor;
uniform vec4 uFogState; // fog_near, 1.0/fog_range

// Lighting environment
uniform vec4 uAmbientColor;

uniform vec4 uLightPos[8]; // pos.xyz, 1.0/radius
uniform vec4 uLightDir[8]; // dir.xyz, inner_rim
uniform vec4 uLightDiffuse[8]; // diffuse.xyz, outer_rim
uniform vec4 uLightSpecular[8]; // specular.xyz, pssm_bias

uniform mat4 uLinearShadowMatrix[4]; // slot 0: linear PSSM shadow matrices
uniform vec4 uLinearShadowSlice; // slot 0: PSSM slice distances linear light
uniform mat4 uSpotShadowMatrix; // slot 1: spot shadow matrix
uniform vec4 uShadowState; // slot 0: inverse resolution, slot1: inverse resolution, slot0: bias, slot1: bias

SAMPLER2DSHADOW(uLinearShadowMap, 14);
SAMPLER2DSHADOW(uSpotShadowMap, 15);

//
mat3 MakeMat3(vec3 c0, vec3 c1, vec3 c2) {
	return mat3(c0, c1, c2);
}

vec3 GetT(mat4 m) { 
#if BGFX_SHADER_LANGUAGE_GLSL
	return vec3(m[3][0], m[3][1], m[3][2]);
#else
	return vec3(m[0][3], m[1][3], m[2][3]);
#endif
}
