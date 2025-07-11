-- ----------------------------
-- Sequence structure for actividades_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."actividades_id_seq";
CREATE SEQUENCE "public"."actividades_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2000000
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for roles_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."roles_id_seq";
CREATE SEQUENCE "public"."roles_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2000000
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for turnos_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."turnos_id_seq";
CREATE SEQUENCE "public"."turnos_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2000000
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for usuarios_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."usuarios_id_seq";
CREATE SEQUENCE "public"."usuarios_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2000000
START 1
CACHE 1;

-- ----------------------------
-- Table structure for actividades
-- ----------------------------
DROP TABLE IF EXISTS "public"."actividades";
CREATE TABLE "public"."actividades" (
  "id" int4 NOT NULL DEFAULT nextval('actividades_id_seq'::regclass),
  "nombre" varchar(150) COLLATE "pg_catalog"."default" NOT NULL,
  "estatus" bool NOT NULL DEFAULT true,
  "creado_en" timestamp(0) NOT NULL DEFAULT now(),
  "actualizado_en" timestamp(0)
)
;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS "public"."roles";
CREATE TABLE "public"."roles" (
  "id" int4 NOT NULL DEFAULT nextval('roles_id_seq'::regclass),
  "nombre" varchar(150) COLLATE "pg_catalog"."default" NOT NULL,
  "estatus" bool NOT NULL DEFAULT true,
  "creado_en" timestamp(0) NOT NULL DEFAULT now(),
  "actualizado_en" timestamp(0)
)
;

-- ----------------------------
-- Table structure for turnos
-- ----------------------------
DROP TABLE IF EXISTS "public"."turnos";
CREATE TABLE "public"."turnos" (
  "id" int4 NOT NULL DEFAULT nextval('turnos_id_seq'::regclass),
  "numero_turno" int4 NOT NULL,
  "actividad_id" int4 NOT NULL,
  "tiempo_recepcion" timestamp(0) NOT NULL DEFAULT now(),
  "tiempo_inicio_atencion" timestamp(0),
  "tiempo_fin_atencion" timestamp(0),
  "usuario_recepcion_id" int4 NOT NULL,
  "usuario_inicio_atencion_id" int4,
  "usuario_fin_atencion_id" int4
)
;

-- ----------------------------
-- Table structure for usuarios
-- ----------------------------
DROP TABLE IF EXISTS "public"."usuarios";
CREATE TABLE "public"."usuarios" (
  "id" int4 NOT NULL DEFAULT nextval('usuarios_id_seq'::regclass),
  "usuario" varchar(25) COLLATE "pg_catalog"."default" NOT NULL,
  "nombre_completo" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "password" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "estatus" bool NOT NULL DEFAULT true,
  "actividad_id" int4 NOT NULL,
  "rol_id" int4 NOT NULL,
  "mesa" int4 NOT NULL,
  "creado_en" timestamp(0) NOT NULL DEFAULT now(),
  "actualizado_en" timestamp(0)
)
;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."actividades_id_seq"
OWNED BY "public"."actividades"."id";

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."roles_id_seq"
OWNED BY "public"."roles"."id";

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."turnos_id_seq"
OWNED BY "public"."turnos"."id";

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."usuarios_id_seq"
OWNED BY "public"."usuarios"."id";

-- ----------------------------
-- Primary Key structure for table actividades
-- ----------------------------
ALTER TABLE "public"."actividades" ADD CONSTRAINT "actividades_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table roles
-- ----------------------------
ALTER TABLE "public"."roles" ADD CONSTRAINT "roles_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table turnos
-- ----------------------------
ALTER TABLE "public"."turnos" ADD CONSTRAINT "turnos_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table usuarios
-- ----------------------------
ALTER TABLE "public"."usuarios" ADD CONSTRAINT "usuarios_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table turnos
-- ----------------------------
ALTER TABLE "public"."turnos" ADD CONSTRAINT "turnos_actividad_id_fkey" FOREIGN KEY ("actividad_id") REFERENCES "public"."actividades" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."turnos" ADD CONSTRAINT "turnos_usuario_fin_atencion_id_fkey" FOREIGN KEY ("usuario_fin_atencion_id") REFERENCES "public"."usuarios" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."turnos" ADD CONSTRAINT "turnos_usuario_inicio_atencion_id_fkey" FOREIGN KEY ("usuario_inicio_atencion_id") REFERENCES "public"."usuarios" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."turnos" ADD CONSTRAINT "turnos_usuario_recepcion_id_fkey" FOREIGN KEY ("usuario_recepcion_id") REFERENCES "public"."usuarios" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table usuarios
-- ----------------------------
ALTER TABLE "public"."usuarios" ADD CONSTRAINT "usuarios_actividad_id_fkey" FOREIGN KEY ("actividad_id") REFERENCES "public"."actividades" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."usuarios" ADD CONSTRAINT "usuarios_rol_id_fkey" FOREIGN KEY ("rol_id") REFERENCES "public"."roles" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
