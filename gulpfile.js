var gulp = require('gulp');
var jshint = require('gulp-jshint');
var concat = require('gulp-concat');
var rename = require('gulp-rename');
var uglify = require('gulp-uglify');
var less = require('gulp-less');
var path = require('path');

// Inspired by https://travismaynard.com/writing/no-need-to-grunt-take-a-gulp-of-fresh-air

// Linting
gulp.task('lint', function() {
  return gulp.src('src/js/**/*.js')
    .pipe(jshint())
    .pipe(jshint.reporter('default'));
});

gulp.task('minify-lib', function() {
  return gulp.src([
      './node_modules/jquery/dist/jquery.min.js',
      './node_modules/underscore/underscore-min.js',
      './node_modules/backbone/backbone-min.js',
    ])
    .pipe(concat('lib.js'))
    .pipe(gulp.dest('./dist/js'));
});

// Concat and minify js
gulp.task('minify', function() {
  return gulp.src('src/js/**/*.js')
    .pipe(concat('app.js'))
    // .pipe(gulp.dest('./dist/js'))
    // .pipe(rename('app.min.js'))
    .pipe(uglify())
    .pipe(gulp.dest('./dist/js'));
});

gulp.task('less-lib', function() {
  return gulp.src('./node_modules/bootstrap/dist/css/bootstrap.min.css')
    .pipe(rename('lib.css'))
    .pipe(gulp.dest('./dist/css'))
});

// Concat and minify less -> css
gulp.task('less', function() {
  return gulp.src('./src/less/**/*.less')
    .pipe(less())
    .pipe(gulp.dest('./dist/css'));
});

gulp.task('build', [
  'less-lib',
  'less',
  'lint',
  'minify-lib',
  'minify',
]);

gulp.task('watch', function() {
  gulp.watch('./src/less/**/*.less', ['less']);
  gulp.watch('./src/js/**/*.js', ['lint', 'minify']);
});

gulp.task('default', ['build', 'watch']);
